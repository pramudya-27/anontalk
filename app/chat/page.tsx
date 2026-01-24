"use client";

import React, {useEffect, useState, useRef} from "react";
import {useRouter} from "next/navigation";
import {apiClient} from "@/lib/api-client";
import {Button} from "@/components/ui/button";
import {Input} from "@/components/ui/input";
import {Card} from "@/components/ui/card";

interface Message {
  id: string;
  senderId: string;
  content: string;
  type?: string;
}

interface ChatSession {
  id: string;
  user1Id: string;
  user2Id: string;
  isActive: boolean;
}

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export default function ChatPage() {
  const router = useRouter();

  const [session, setSession] = useState<ChatSession | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(true);
  const [currentUser, setCurrentUser] = useState<string | null>(null);
  const [isSearching, setIsSearching] = useState(false);
  const [isUploading, setIsUploading] = useState(false);

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    const initializeChat = async () => {
      try {
        const user = await apiClient.getMe();
        if (!user) {
          router.push("/auth/login");
          return;
        }

        setCurrentUser(user.id);

        // Get active session
        const sessionData = await apiClient.getSessions();
        if (sessionData.sessions && sessionData.sessions.length > 0) {
          const activeSession = sessionData.sessions[0];
          setSession(activeSession);
          fetchMessages(activeSession.id);
        } else {
          // If no active session, start searching for one
          setIsSearching(true);
          try {
            const newSession = await apiClient.findMatch();
            setSession(newSession);
            fetchMessages(newSession.id);
          } catch (e) {
            console.log("No match found yet");
          }
          setIsSearching(false);
        }

        setLoading(false);
      } catch (error) {
        console.error("Error initializing chat:", error);
        setLoading(false);
      }
    };

    initializeChat();
  }, [router]);

  useEffect(() => {
    if (!session) return;

    console.log("Connecting to WebSocket for session:", session.id);
    const ws = apiClient.connectWebSocket(
      session.id,
      (message: any) => {
        console.log("Received WebSocket message:", message);
        if (message.type === "session_ended") {
          alert("The chat has been ended.");
          router.push("/dashboard");
        } else {
          setMessages((prev) => {
            // Avoid duplicates if any
            if (prev.some((m) => m.id === message.id)) return prev;
            return [...prev, message];
          });
        }
      },
      (error) => {
        console.error("WebSocket error:", error);
      },
    );

    return () => {
      console.log("Closing WebSocket connection");
      ws.close();
    };
  }, [session, router]);

  const fetchMessages = async (sessionId: string) => {
    try {
      const data = await apiClient.getMessages(sessionId);
      if (data && data.messages) {
        setMessages(data.messages);
        scrollToBottom();
      }
    } catch (error) {
      console.error("Error fetching messages:", error);
    }
  };

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({behavior: "smooth"});
    // Also scroll window bottom just in case
    window.scrollTo({top: document.body.scrollHeight, behavior: "smooth"});
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSendMessage = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!input.trim() || !session || !currentUser) return;

    try {
      await apiClient.sendMessage(session.id, input, "text");
      setInput("");
      // Message will be received via WebSocket
    } catch (error) {
      console.error("Error sending message:", error);
    }
  };

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file || !session) return;

    setIsUploading(true);
    try {
      // 1. Upload file
      const {url} = await apiClient.uploadFile(session.id, file);

      // 2. Determine type
      let type = "file";
      if (file.type.startsWith("image/")) {
        type = "image";
      } else if (file.type.startsWith("audio/")) {
        type = "audio";
      }

      // 3. Send message with URL and type
      await apiClient.sendMessage(session.id, url, type);
    } catch (error) {
      console.error("Error uploading file:", error);
      alert("Failed to upload file");
    } finally {
      setIsUploading(false);
      // Reset input
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
    }
  };

  const handleEndChat = async () => {
    if (!session) return;

    try {
      await apiClient.endSession(session.id);
      router.push("/dashboard");
    } catch (error) {
      console.error("Error ending session:", error);
    }
  };

  const handleReport = async () => {
    if (!session || !currentUser) return;

    const otherUserId =
      session.user1Id === currentUser ? session.user2Id : session.user1Id;

    const reason = prompt("Please describe the reason for reporting:");
    if (!reason) return;

    try {
      await apiClient.reportUser(
        otherUserId,
        reason,
        "User reported from chat",
      );
      alert("Report submitted. Thank you for keeping AnonTalk safe.");
      handleEndChat();
    } catch (error) {
      console.error("Error reporting user:", error);
      alert("Failed to report user.");
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <p className="text-lg">Loading chat...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      {/* Header */}
      <header className="bg-white shadow-sm py-4 px-6 sticky top-0 z-10 flex justify-between items-center">
        <div>
          <h1 className="text-xl font-bold text-gray-800">Anonymous Chat</h1>
          <p className="text-xs text-gray-500">Matched with a stranger</p>
        </div>
        <div className="flex gap-2">
          <button
            className="px-3 py-1 text-sm text-red-600 hover:bg-red-50 rounded border border-red-200"
            onClick={handleReport}
          >
            Report
          </button>
          <button
            onClick={handleEndChat}
            className="px-3 py-1 text-sm bg-white border border-gray-300 text-gray-700 rounded hover:bg-gray-50"
          >
            End Chat
          </button>
        </div>
      </header>

      {/* Chat Area */}
      <div className="flex-1 p-4 pb-24 overflow-y-auto min-h-0">
        <div className="max-w-3xl mx-auto space-y-4">
          {isSearching ? (
            <div className="flex flex-col items-center justify-center py-20 text-gray-500">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-500 mb-4"></div>
              <p>Looking for a match...</p>
            </div>
          ) : session ? (
            <>
              {messages.map((msg) => {
                const isMe = msg.senderId === currentUser;
                return (
                  <div
                    key={msg.id}
                    className={`flex ${isMe ? "justify-end" : "justify-start"}`}
                  >
                    <div
                      className={`max-w-[70%] rounded-2xl px-4 py-2 ${
                        isMe
                          ? "bg-blue-600 text-white rounded-br-none"
                          : "bg-white text-gray-800 rounded-bl-none shadow-sm"
                      }`}
                    >
                      {msg.type === "image" ? (
                        <div className="relative group">
                          <img
                            src={`${API_BASE_URL}${msg.content}`}
                            alt="Sent image"
                            className="max-w-full rounded-lg cursor-pointer hover:opacity-90"
                            onClick={() =>
                              window.open(
                                `${API_BASE_URL}${msg.content}`,
                                "_blank",
                              )
                            }
                          />
                          <a
                            href={`${API_BASE_URL}${msg.content}`}
                            download
                            target="_blank"
                            rel="noopener noreferrer"
                            className="absolute bottom-2 right-2 p-1.5 bg-black/50 text-white rounded-full opacity-0 group-hover:opacity-100 transition-opacity hover:bg-black/70"
                            title="Download Image"
                            onClick={(e) => e.stopPropagation()}
                          >
                            <svg
                              xmlns="http://www.w3.org/2000/svg"
                              fill="none"
                              viewBox="0 0 24 24"
                              strokeWidth={2}
                              stroke="currentColor"
                              className="w-4 h-4"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M12 12.75l-3-3m3 3l3-3m-3 3V3"
                              />
                            </svg>
                          </a>
                        </div>
                      ) : msg.type === "audio" ? (
                        <div className="flex flex-col gap-1 min-w-[240px]">
                          <audio controls className="w-full">
                            <source src={`${API_BASE_URL}${msg.content}`} />
                            Your browser does not support the audio element.
                          </audio>
                          <div className="flex justify-end">
                            <a
                              href={`${API_BASE_URL}${msg.content}`}
                              download
                              target="_blank"
                              className={`text-xs flex items-center gap-1 hover:underline ${isMe ? "text-blue-100" : "text-blue-600"}`}
                            >
                              <svg
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 24 24"
                                strokeWidth={1.5}
                                stroke="currentColor"
                                className="w-3 h-3"
                              >
                                <path
                                  strokeLinecap="round"
                                  strokeLinejoin="round"
                                  d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M12 12.75l-3-3m3 3l3-3m-3 3V3"
                                />
                              </svg>
                              Download Audio
                            </a>
                          </div>
                        </div>
                      ) : msg.type === "file" ? (
                        <a
                          href={`${API_BASE_URL}${msg.content}`}
                          target="_blank"
                          rel="noreferrer"
                          download
                          className={`flex items-center gap-2 underline decoration-dotted underline-offset-4 py-1 px-2 hover:bg-black/5 rounded transition-colors ${isMe ? "text-white hover:bg-white/10" : "text-blue-600"}`}
                        >
                          <span className="text-xl">📄</span>
                          <span>Open/Download File</span>
                        </a>
                      ) : (
                        <p className="break-words white-space-pre-wrap">
                          {msg.content}
                        </p>
                      )}
                    </div>
                  </div>
                );
              })}
              <div ref={messagesEndRef} />
            </>
          ) : (
            <div className="flex flex-col items-center justify-center py-20 text-gray-500 hover:text-gray-700 transition">
              <p className="mb-4">No active chat found.</p>
              <Button onClick={() => window.location.reload()}>
                Retry Match
              </Button>
            </div>
          )}
        </div>
      </div>

      {/* Input Area */}
      {session && (
        <div className="fixed bottom-0 left-0 right-0 bg-white border-t p-4">
          <div className="max-w-3xl mx-auto">
            <form
              onSubmit={handleSendMessage}
              className="flex gap-2 items-center"
            >
              <input
                type="file"
                ref={fileInputRef}
                className="hidden"
                onChange={handleFileUpload}
              />
              <button
                type="button"
                onClick={() => fileInputRef.current?.click()}
                disabled={isUploading}
                className="p-2 text-gray-500 hover:bg-gray-100 rounded-full transition-colors flex-shrink-0"
                title="Attach file"
              >
                {isUploading ? (
                  <div className="animate-spin h-5 w-5 border-2 border-gray-400 border-t-transparent rounded-full" />
                ) : (
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="w-6 h-6"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M18.375 12.739l-7.693 7.693a4.5 4.5 0 01-6.364-6.364l10.94-10.94A3 3 0 1119.5 6.375L8.545 17.365a1.5 1.5 0 01-2.121-2.121l10.859-10.86"
                    />
                  </svg>
                )}
              </button>

              <div className="flex-1">
                <Input
                  type="text"
                  value={input}
                  onChange={(e) => setInput(e.target.value)}
                  placeholder="Type a message..."
                  className="w-full rounded-full"
                />
              </div>

              <button
                type="submit"
                disabled={!input.trim()}
                className="bg-gray-800 text-white px-6 py-2 rounded-full hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed font-medium transition-colors flex-shrink-0"
              >
                Send
              </button>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
