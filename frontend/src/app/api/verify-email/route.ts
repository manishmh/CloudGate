import { NextRequest, NextResponse } from 'next/server';

export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url);
    const token = searchParams.get('token');
    const userId = searchParams.get('userId');

    if (!token || !userId) {
      return NextResponse.redirect(
        new URL('/profile?verification=invalid', request.url)
      );
    }

    // In a real app, you'd verify the token against your database
    // For demo purposes, we'll simulate verification
    console.log('Verifying email for user:', userId, 'with token:', token);

    // Simulate token validation (in production, check database)
    // For demo, we'll assume the token is valid if it exists
    if (token.length === 64) { // Our tokens are 32 bytes = 64 hex chars
      // Mark email as verified in your database
      // For demo, we'll just redirect with success
      
      return NextResponse.redirect(
        new URL('/profile?verification=success', request.url)
      );
    } else {
      return NextResponse.redirect(
        new URL('/profile?verification=invalid', request.url)
      );
    }

  } catch (error) {
    console.error('Error verifying email:', error);
    return NextResponse.redirect(
      new URL('/profile?verification=error', request.url)
    );
  }
} 