import type { Metadata } from "next";
import Link from "next/link";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { NotificationBell } from "@/components/notification-bell";
import { SiteFooter } from "@/components/site-footer";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "FoodLike",
  description: "みんなの好き嫌いにこっそり配慮して外食先を提案するアプリ",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja" className={`${geistSans.variable} ${geistMono.variable}`}>
      <body>
        <header className="header">
          <div className="header-inner">
            <Link href="/" className="logo">
              🍴 FoodLike
            </Link>
            <NotificationBell />
          </div>
        </header>
        <main className="container">{children}</main>
        <SiteFooter />
      </body>
    </html>
  );
}
