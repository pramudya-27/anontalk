const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export class APIClient {
  private token: string | null = null;

  constructor() {
    if (typeof window !== "undefined") {
      this.token = localStorage.getItem("token");
    }
  }

  setToken(token: string) {
    this.token = token;
    localStorage.setItem("token", token);
  }

  getToken() {
    return this.token;
  }

  clearToken() {
    this.token = null;
    localStorage.removeItem("token");
  }

  private async request(method: string, endpoint: string, body?: any) {
    const headers: any = {
      "Content-Type": "application/json",
    };

    if (this.token) {
      headers["Authorization"] = `Bearer ${this.token}`;
    }

    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });

    if (!response.ok) {
      const error = await response.text();
      throw new Error(error || `API Error: ${response.status}`);
    }

    return response.json();
  }

  async get(endpoint: string) {
    return this.request("GET", endpoint);
  }

  async post(endpoint: string, body: any) {
    return this.request("POST", endpoint, body);
  }

  async put(endpoint: string, body: any) {
    return this.request("PUT", endpoint, body);
  }

  // Auth endpoints
  async signup(
    email: string,
    password: string,
    gender: string,
    preferredGender: string[],
    interests: string[],
  ) {
    const response = await this.post("/api/auth/signup", {
      email,
      password,
      gender,
      preferredGender,
      interests,
    });
    this.setToken(response.token);
    return response;
  }

  async login(email: string, password: string) {
    const response = await this.post("/api/auth/login", {
      email,
      password,
    });
    this.setToken(response.token);
    return response;
  }

  async logout() {
    try {
      await this.post("/api/auth/logout", {});
    } catch (error) {
      // Ignore errors during logout (e.g. if token is already invalid)
      console.error("Logout error:", error);
    } finally {
      this.clearToken();
    }
  }

  async getMe() {
    if (!this.token) return null;
    return this.get("/api/auth/me");
  }

  // User endpoints
  async getProfile() {
    return this.get("/api/users/profile");
  }

  async updateProfile(
    displayName: string,
    ageRange: string,
    lookingFor: string,
  ) {
    return this.put("/api/users/profile", {
      displayName,
      ageRange,
      lookingFor,
    });
  }

  async updatePreferences(preferredGender: string[], interests: string[]) {
    return this.put("/api/users/preferences", {
      preferredGender,
      interests,
    });
  }

  // Matchmaking endpoints
  async joinQueue(
    gender: string,
    preferredGender: string[],
    interests: string[],
  ) {
    return this.post("/api/matchmaking/join-queue", {
      gender,
      preferredGender,
      interests,
    });
  }

  async leaveQueue() {
    return this.post("/api/matchmaking/leave-queue", {});
  }

  async getMatchStatus() {
    return this.get("/api/matchmaking/status");
  }

  async findMatch() {
    // simplified: join queue and poll for status
    // For this demo, we assume the backend might match immediately or we just join
    await this.joinQueue("any", [], []);

    // Poll for a valid session for a few seconds
    for (let i = 0; i < 10; i++) {
      await new Promise((r) => setTimeout(r, 1000));
      const sessionRes = await this.getSessions();
      if (sessionRes.sessions && sessionRes.sessions.length > 0) {
        return sessionRes.sessions[0];
      }
    }
    throw new Error("No match found");
  }

  // Chat endpoints
  async getSessions() {
    return this.get("/api/chat/sessions");
  }

  async getMessages(sessionId: string) {
    return this.get(`/api/chat/sessions/${sessionId}/messages`);
  }

  async sendMessage(sessionId: string, content: string, type: string = "text") {
    return this.post("/api/chat/send-message", {
      sessionId,
      content,
      type,
    });
  }

  async uploadFile(sessionId: string, file: File) {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("sessionId", sessionId);

    // We use raw fetch here because our request wrapper sets Content-Type to application/json
    const headers: any = {};
    if (this.token) {
      headers["Authorization"] = `Bearer ${this.token}`;
    }

    const response = await fetch(`${API_BASE_URL}/api/chat/upload`, {
      method: "POST",
      headers,
      body: formData,
    });

    if (!response.ok) {
      const error = await response.text();
      throw new Error(error || `API Error: ${response.status}`);
    }

    return response.json();
  }

  async endSession(sessionId: string) {
    return this.post(`/api/chat/sessions/${sessionId}/end`, {});
  }

  // Report endpoints
  async reportUser(
    reportedUserId: string,
    reason: string,
    description: string,
  ) {
    return this.post("/api/reports/create", {
      reportedUserId,
      reason,
      description,
    });
  }

  // Admin endpoints
  async getStats() {
    return this.get("/api/admin/stats");
  }

  async getReports() {
    return this.get("/api/admin/reports");
  }

  async handleReport(reportId: string, action: string, adminNotes: string) {
    return this.post(`/api/admin/reports/${reportId}/action`, {
      action,
      adminNotes,
    });
  }

  async banUser(userId: string, reason: string) {
    return this.post(`/api/admin/users/${userId}/ban`, {
      reason,
    });
  }

  // WebSocket helper
  connectWebSocket(
    sessionId: string,
    onMessage: (message: any) => void,
    onError: (error: any) => void,
  ) {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = `${protocol}//${window.location.host.split(":")[0]}:8080/ws/chat/${sessionId}?token=${this.token}`;

    const ws = new WebSocket(wsUrl);

    ws.onmessage = (event) => {
      onMessage(JSON.parse(event.data));
    };

    ws.onerror = (error) => {
      onError(error);
    };

    return ws;
  }
}

export const apiClient = new APIClient();
