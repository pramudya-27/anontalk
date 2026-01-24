"use client";

import React from "react";

import {useState} from "react";
import {apiClient} from "@/lib/api-client";
import {useRouter} from "next/navigation";
import {Button} from "@/components/ui/button";
import {Input} from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export default function SetupProfilePage() {
  const [username, setUsername] = useState("");
  const [gender, setGender] = useState<"male" | "female" | "other" | "">("");
  const [interestedIn, setInterestedIn] = useState<
    "male" | "female" | "both" | ""
  >("");
  const [interests, setInterests] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    if (!username || !gender || !interestedIn) {
      setError("Please fill in all fields");
      setLoading(false);
      return;
    }

    try {
      // Update profile (display name, age range, looking for)
      await apiClient.updateProfile(username, "18-99", interestedIn);

      // Update preferences based on interestedIn
      let preferredGender: string[] = [];
      if (interestedIn === "male") preferredGender = ["male"];
      else if (interestedIn === "female") preferredGender = ["female"];
      else if (interestedIn === "both") preferredGender = ["male", "female"];

      const interestsList = interests
        .split(",")
        .map((i) => i.trim())
        .filter((i) => i.length > 0);

      await apiClient.updatePreferences(preferredGender, interestsList);

      router.push("/dashboard"); // Redirect to dashboard instead of chat directly
    } catch (err: any) {
      setError(err.message || "Failed to update profile");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-background p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-2">
          <CardTitle className="text-2xl">Complete Your Profile</CardTitle>
          <CardDescription>
            Set up your preferences for matching
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            {error && (
              <div className="bg-destructive/10 text-destructive text-sm p-3 rounded-md">
                {error}
              </div>
            )}
            <div className="space-y-2">
              <label htmlFor="username" className="text-sm font-medium">
                Username
              </label>
              <Input
                id="username"
                placeholder="Choose a username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                required
              />
            </div>
            <div className="space-y-2">
              <label htmlFor="gender" className="text-sm font-medium">
                Your Gender
              </label>
              <select
                id="gender"
                value={gender}
                onChange={(e) => setGender(e.target.value as any)}
                className="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground"
                required
              >
                <option value="">Select your gender</option>
                <option value="male">Male</option>
                <option value="female">Female</option>
                <option value="other">Other</option>
              </select>
            </div>
            <div className="space-y-2">
              <label htmlFor="interestedIn" className="text-sm font-medium">
                Interested In
              </label>
              <select
                id="interestedIn"
                value={interestedIn}
                onChange={(e) => setInterestedIn(e.target.value as any)}
                className="w-full px-3 py-2 border border-input rounded-md bg-background text-foreground"
                required
              >
                <option value="">Who do you want to chat with?</option>
                <option value="male">Males</option>
                <option value="female">Females</option>
                <option value="both">Both</option>
              </select>
            </div>
            <div className="space-y-2">
              <label htmlFor="interests" className="text-sm font-medium">
                Interests (comma separated)
              </label>
              <Input
                id="interests"
                placeholder="e.g. Music, Tech, Travel"
                value={interests}
                onChange={(e) => setInterests(e.target.value)}
              />
            </div>
            <Button type="submit" disabled={loading} className="w-full">
              {loading ? "Setting up..." : "Continue to Dashboard"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
