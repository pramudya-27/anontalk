"use client";

import {useEffect, useState} from "react";
import {useRouter} from "next/navigation";
import {apiClient} from "@/lib/api-client";
import {Button} from "@/components/ui/button";
import {Card} from "@/components/ui/card";
import {Heart, Users, Shield, LogOut} from "lucide-react";

interface Profile {
  id: string;
  username: string; // mapped from displayName
  gender: string;
  seeking_gender: string; // mapped from lookingFor
  interests: string[];
}

export default function DashboardPage() {
  const router = useRouter();

  const [profile, setProfile] = useState<Profile | null>(null);
  const [loading, setLoading] = useState(true);
  const [isMatchmaking, setIsMatchmaking] = useState(false);
  const [matchmakingStatus, setMatchmakingStatus] = useState("");

  useEffect(() => {
    const loadProfile = async () => {
      try {
        const user = await apiClient.getMe();
        if (!user) {
          router.push("/auth/login");
          return;
        }

        // Get full profile data
        const profileData = await apiClient.getProfile();

        // Check if profile is complete (needs displayName)
        if (!profileData.displayName) {
          router.push("/setup-profile");
          return;
        }

        setProfile({
          id: profileData.id,
          username: profileData.displayName,
          gender: profileData.gender,
          seeking_gender: profileData.lookingFor,
          interests: profileData.interests || [],
        });
        setLoading(false);
      } catch (error) {
        console.error("Error loading profile:", error);
        router.push("/auth/login");
      }
    };

    loadProfile();
  }, []);

  const startMatchmaking = async () => {
    if (!profile) return;

    setIsMatchmaking(true);
    setMatchmakingStatus("Finding a match...");

    try {
      // Join matchmaking queue
      // Note: apiClient.joinQueue expects arrays for preferredGender
      let preferredGender: string[] = [];
      if (profile.seeking_gender === "male") preferredGender = ["male"];
      else if (profile.seeking_gender === "female")
        preferredGender = ["female"];
      else if (profile.seeking_gender === "both")
        preferredGender = ["male", "female"];

      const response = await apiClient.joinQueue(
        profile.gender,
        preferredGender,
        profile.interests,
      );

      if (response.match) {
        // Match found immediately
        setMatchmakingStatus("Match found! Starting chat...");
        setTimeout(() => {
          router.push("/chat");
        }, 1500);
      } else {
        setMatchmakingStatus("Waiting for a match...");
        // Poll for match
        let pollCount = 0;
        const pollInterval = setInterval(async () => {
          pollCount++;

          try {
            const status = await apiClient.getMatchStatus();
            if (status.match) {
              clearInterval(pollInterval);
              setMatchmakingStatus("Match found! Starting chat...");
              setTimeout(() => {
                router.push("/chat");
              }, 1500);
              return;
            }
          } catch (e) {
            console.error("Error polling match status", e);
          }

          if (pollCount > 60) {
            clearInterval(pollInterval);
            setIsMatchmaking(false);
            setMatchmakingStatus("No match found. Please try again later.");
            // Optionally leave queue
            await apiClient.leaveQueue();
          }
        }, 2000);
      }
    } catch (error) {
      console.error("Matchmaking error:", error);
      setIsMatchmaking(false);
      setMatchmakingStatus("Error starting matchmaking. Please try again.");
    }
  };

  const handleLogout = async () => {
    await apiClient.logout();
    router.push("/");
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <p className="text-lg">Loading profile...</p>
      </div>
    );
  }

  return (
    <main className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      {/* Header */}
      <div className="bg-white border-b shadow-sm">
        <div className="max-w-6xl mx-auto px-4 py-4 flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">DiAnonTalk</h1>
            <p className="text-sm text-gray-500">Anonymous Chat Matching</p>
          </div>
          <Button
            variant="outline"
            size="sm"
            onClick={handleLogout}
            className="gap-2 bg-transparent"
          >
            <LogOut className="w-4 h-4" />
            Logout
          </Button>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-6xl mx-auto px-4 py-12">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Profile Card */}
          <Card className="lg:col-span-1 p-6">
            <div className="space-y-4">
              <div>
                <h2 className="text-lg font-semibold text-gray-900">
                  {profile?.username}
                </h2>
                <p className="text-sm text-gray-500">
                  {profile?.gender} seeking {profile?.seeking_gender}
                </p>
              </div>

              <div className="pt-4 border-t">
                <h3 className="font-medium text-gray-900 mb-2">Interests</h3>
                <div className="flex flex-wrap gap-2">
                  {profile?.interests?.map((interest) => (
                    <span
                      key={interest}
                      className="inline-block bg-blue-100 text-blue-800 text-xs px-3 py-1 rounded-full"
                    >
                      {interest}
                    </span>
                  ))}
                  {(!profile?.interests || profile.interests.length === 0) && (
                    <span className="text-sm text-gray-500">
                      No interests set
                    </span>
                  )}
                </div>
              </div>

              <Button
                variant="outline"
                className="w-full mt-4 bg-transparent"
                onClick={() => router.push("/setup-profile")}
              >
                Edit Profile
              </Button>
            </div>
          </Card>

          {/* Matchmaking Section */}
          <div className="lg:col-span-2 space-y-6">
            {/* Start Matching Card */}
            <Card className="p-8 text-center bg-gradient-to-r from-blue-500 to-indigo-600 text-white">
              <Users className="w-12 h-12 mx-auto mb-4" />
              <h2 className="text-2xl font-bold mb-2">Ready to Chat?</h2>
              <p className="mb-6 text-blue-100">
                Connect with someone new who matches your preferences
              </p>
              <Button
                size="lg"
                disabled={isMatchmaking}
                onClick={startMatchmaking}
                className="bg-white text-blue-600 hover:bg-blue-50"
              >
                {isMatchmaking ? (
                  <>
                    <span className="inline-block animate-spin mr-2">⏳</span>
                    {matchmakingStatus || "Matching..."}
                  </>
                ) : (
                  <>
                    <Heart className="w-5 h-5 mr-2" />
                    Start Matchmaking
                  </>
                )}
              </Button>
            </Card>

            {/* Features */}
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <Card className="p-4">
                <Shield className="w-8 h-8 text-green-600 mb-2" />
                <h3 className="font-semibold mb-1">Complete Anonymity</h3>
                <p className="text-sm text-gray-600">
                  Your real identity is always hidden
                </p>
              </Card>
              <Card className="p-4">
                <Users className="w-8 h-8 text-blue-600 mb-2" />
                <h3 className="font-semibold mb-1">Smart Matching</h3>
                <p className="text-sm text-gray-600">
                  Connected based on your preferences
                </p>
              </Card>
            </div>

            {/* Safety Info */}
            <Card className="p-4 bg-yellow-50 border-yellow-200">
              <h3 className="font-semibold text-yellow-900 mb-2">
                Safety Reminder
              </h3>
              <ul className="text-sm text-yellow-800 space-y-1">
                <li>• Never share personal identifying information</li>
                <li>• Report any inappropriate behavior immediately</li>
                <li>• Keep conversations respectful and appropriate</li>
              </ul>
            </Card>
          </div>
        </div>
      </div>
    </main>
  );
}
