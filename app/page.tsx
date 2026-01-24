"use client";

import {useEffect, useState} from "react";
import {useRouter} from "next/navigation";
import {apiClient} from "@/lib/api-client";
import {Button} from "@/components/ui/button";
import {Card} from "@/components/ui/card";
import {Heart, Shield, Zap, Users, ArrowRight} from "lucide-react";

export default function Home() {
  const router = useRouter();
  const [user, setUser] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const checkUser = async () => {
      try {
        const user = await apiClient.getMe();
        setUser(user);
      } catch (error) {
        // Not logged in
      } finally {
        setLoading(false);
      }
    };

    checkUser();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        Loading...
      </div>
    );
  }

  if (user) {
    router.push("/dashboard");
    return null;
  }

  return (
    <main className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50">
      {/* Navigation */}
      <nav className="fixed top-0 w-full bg-white/80 backdrop-blur-md border-b z-50">
        <div className="max-w-6xl mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Heart className="w-6 h-6 text-blue-600" />
            <h1 className="text-xl font-bold text-gray-900">DiAnonTalk</h1>
          </div>
          <div className="flex gap-3">
            <Button
              variant="outline"
              onClick={() => router.push("/auth/login")}
            >
              Sign In
            </Button>
            <Button onClick={() => router.push("/auth/signup")}>
              Get Started
            </Button>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="pt-32 pb-20 px-4">
        <div className="max-w-4xl mx-auto text-center">
          <h2 className="text-5xl md:text-6xl font-bold text-gray-900 mb-6">
            Connect Anonymously,
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-indigo-600">
              {" "}
              Naturally
            </span>
          </h2>
          <p className="text-xl text-gray-600 mb-8 max-w-2xl mx-auto">
            Meet new people with matching interests in completely anonymous
            chats. No profiles, no pressure, just genuine conversations.
          </p>
          <div className="flex gap-4 justify-center">
            <Button size="lg" onClick={() => router.push("/auth/signup")}>
              Start Chatting <ArrowRight className="w-5 h-5 ml-2" />
            </Button>
            <Button
              size="lg"
              variant="outline"
              onClick={() => {
                const element = document.getElementById("features");
                element?.scrollIntoView({behavior: "smooth"});
              }}
            >
              Learn More
            </Button>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="py-20 px-4 bg-white">
        <div className="max-w-6xl mx-auto">
          <h3 className="text-3xl font-bold text-center mb-16">
            Why AnonTalk?
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <Card className="p-8 hover:shadow-lg transition-shadow">
              <Shield className="w-12 h-12 text-blue-600 mb-4" />
              <h4 className="text-xl font-semibold mb-2">Complete Anonymity</h4>
              <p className="text-gray-600">
                Your real identity stays private. Connect without revealing who
                you are.
              </p>
            </Card>
            <Card className="p-8 hover:shadow-lg transition-shadow">
              <Users className="w-12 h-12 text-blue-600 mb-4" />
              <h4 className="text-xl font-semibold mb-2">Smart Matching</h4>
              <p className="text-gray-600">
                Find people who match your preferences and interests instantly.
              </p>
            </Card>
            <Card className="p-8 hover:shadow-lg transition-shadow">
              <Zap className="w-12 h-12 text-blue-600 mb-4" />
              <h4 className="text-xl font-semibold mb-2">Real Connections</h4>
              <p className="text-gray-600">
                Have genuine conversations with people who share your passions.
              </p>
            </Card>
          </div>
        </div>
      </section>

      {/* How It Works */}
      <section className="py-20 px-4">
        <div className="max-w-4xl mx-auto">
          <h3 className="text-3xl font-bold text-center mb-16">How It Works</h3>
          <div className="space-y-8">
            <div className="flex gap-6 items-start">
              <div className="bg-blue-600 text-white w-12 h-12 rounded-full flex items-center justify-center font-bold flex-shrink-0">
                1
              </div>
              <div>
                <h4 className="text-xl font-semibold mb-2">Create Account</h4>
                <p className="text-gray-600">
                  Sign up with your email and set your preferences.
                </p>
              </div>
            </div>
            <div className="flex gap-6 items-start">
              <div className="bg-blue-600 text-white w-12 h-12 rounded-full flex items-center justify-center font-bold flex-shrink-0">
                2
              </div>
              <div>
                <h4 className="text-xl font-semibold mb-2">Set Preferences</h4>
                <p className="text-gray-600">
                  Tell us who you'd like to connect with based on interests and
                  gender.
                </p>
              </div>
            </div>
            <div className="flex gap-6 items-start">
              <div className="bg-blue-600 text-white w-12 h-12 rounded-full flex items-center justify-center font-bold flex-shrink-0">
                3
              </div>
              <div>
                <h4 className="text-xl font-semibold mb-2">Get Matched</h4>
                <p className="text-gray-600">
                  Our algorithm finds someone who matches your preferences.
                </p>
              </div>
            </div>
            <div className="flex gap-6 items-start">
              <div className="bg-blue-600 text-white w-12 h-12 rounded-full flex items-center justify-center font-bold flex-shrink-0">
                4
              </div>
              <div>
                <h4 className="text-xl font-semibold mb-2">Start Chatting</h4>
                <p className="text-gray-600">
                  Connect and have genuine conversations with your match.
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 px-4 bg-gradient-to-r from-blue-600 to-indigo-600">
        <div className="max-w-4xl mx-auto text-center text-white">
          <h3 className="text-4xl font-bold mb-6">
            Ready to Meet Someone New?
          </h3>
          <p className="text-lg mb-8 opacity-90">
            Join thousands of people connecting anonymously every day.
          </p>
          <Button
            size="lg"
            onClick={() => router.push("/auth/signup")}
            className="bg-white text-blue-600 hover:bg-blue-50"
          >
            Get Started Now <ArrowRight className="w-5 h-5 ml-2" />
          </Button>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-gray-900 text-gray-400 py-8 px-4">
        <div className="max-w-6xl mx-auto text-center">
          <p>
            &copy; 2025 AnonTalk. All rights reserved. Chat safely and
            respectfully.
          </p>
        </div>
      </footer>
    </main>
  );
}
