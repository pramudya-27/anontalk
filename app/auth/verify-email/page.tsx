import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import Link from 'next/link'
import { Button } from '@/components/ui/button'

export default function VerifyEmailPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-background p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-2">
          <CardTitle className="text-2xl">Verify Your Email</CardTitle>
          <CardDescription>
            We've sent you a confirmation link
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <p className="text-sm text-muted-foreground">
            Check your email for a confirmation link. Click it to verify your account and get started with AnonTalk.
          </p>
          <p className="text-sm text-muted-foreground">
            Didn't receive an email? Check your spam folder or{' '}
            <Link href="/auth/signup" className="text-primary hover:underline">
              try signing up again
            </Link>
            .
          </p>
          <Link href="/auth/login">
            <Button variant="outline" className="w-full bg-transparent">
              Back to Login
            </Button>
          </Link>
        </CardContent>
      </Card>
    </div>
  )
}
