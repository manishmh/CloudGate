import crypto from 'crypto';
import { NextRequest, NextResponse } from 'next/server';
import nodemailer from 'nodemailer';

// Create transporter with the provided credentials
const transporter = nodemailer.createTransport({
  service: 'gmail',
  auth: {
    user: process.env.EMAIL_USER,
    pass: process.env.EMAIL_PASS,
  },
});

export async function POST(request: NextRequest) {
  try {
    const { email, userId, name } = await request.json();

    if (!email || !userId || !name) {
      return NextResponse.json(
        { error: 'Missing required fields' },
        { status: 400 }
      );
    }

    // Generate verification token
    const verificationToken = crypto.randomBytes(32).toString('hex');
    
    // In a real app, you'd store this token in your database with expiration
    // For now, we'll create a simple verification URL
    const verificationUrl = `${process.env.NEXT_PUBLIC_APP_URL || 'http://localhost:3000'}/api/verify-email?token=${verificationToken}&userId=${userId}`;

    // Email template
    const mailOptions = {
      from: process.env.EMAIL_USER,
      to: email,
      subject: 'CloudGate - Verify Your Email Address',
      html: `
        <!DOCTYPE html>
        <html>
        <head>
          <meta charset="utf-8">
          <meta name="viewport" content="width=device-width, initial-scale=1.0">
          <title>Verify Your Email - CloudGate</title>
          <style>
            body {
              font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
              line-height: 1.6;
              color: #333;
              max-width: 600px;
              margin: 0 auto;
              padding: 20px;
            }
            .header {
              background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
              color: white;
              padding: 30px;
              text-align: center;
              border-radius: 10px 10px 0 0;
            }
            .content {
              background: #f8f9fa;
              padding: 30px;
              border-radius: 0 0 10px 10px;
            }
            .button {
              display: inline-block;
              background: #2563eb;
              color: white;
              padding: 12px 30px;
              text-decoration: none;
              border-radius: 6px;
              font-weight: 600;
              margin: 20px 0;
            }
            .footer {
              text-align: center;
              margin-top: 30px;
              padding-top: 20px;
              border-top: 1px solid #e5e7eb;
              color: #6b7280;
              font-size: 14px;
            }
          </style>
        </head>
        <body>
          <div class="header">
            <h1>🔐 CloudGate</h1>
            <p>Secure Single Sign-On Platform</p>
          </div>
          
          <div class="content">
            <h2>Hello ${name}!</h2>
            
            <p>Thank you for using CloudGate. To complete your account setup and ensure the security of your account, please verify your email address.</p>
            
            <p>Click the button below to verify your email:</p>
            
            <div style="text-align: center;">
              <a href="${verificationUrl}" class="button">Verify Email Address</a>
            </div>
            
            <p>If the button doesn't work, you can copy and paste this link into your browser:</p>
            <p style="word-break: break-all; background: #e5e7eb; padding: 10px; border-radius: 4px; font-family: monospace;">
              ${verificationUrl}
            </p>
            
            <p><strong>Important:</strong> This verification link will expire in 24 hours for security reasons.</p>
            
            <p>If you didn't request this verification, please ignore this email or contact our support team.</p>
          </div>
          
          <div class="footer">
            <p>This email was sent by CloudGate SSO Platform</p>
            <p>If you have any questions, please contact our support team.</p>
          </div>
        </body>
        </html>
      `,
    };

    // Send email
    await transporter.sendMail(mailOptions);

    // In a real app, store the verification token in database with expiration
    // For demo purposes, we'll store it in a simple way (not recommended for production)
    const verificationData = {
      token: verificationToken,
      userId,
      email,
      createdAt: new Date().toISOString(),
      expiresAt: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString(), // 24 hours
    };

    // Store in localStorage on server side (not ideal, but for demo)
    // In production, use a proper database
    console.log('Verification token generated:', verificationData);

    return NextResponse.json(
      { 
        message: 'Verification email sent successfully',
        success: true 
      },
      { status: 200 }
    );

  } catch (error) {
    console.error('Error sending verification email:', error);
    return NextResponse.json(
      { error: 'Failed to send verification email' },
      { status: 500 }
    );
  }
} 