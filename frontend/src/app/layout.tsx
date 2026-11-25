import type { Metadata } from "next";
import { Plus_Jakarta_Sans } from "next/font/google";
import "./globals.css";
import { ThemeProvider } from "@/providers/theme-provider";
import { UserProvider } from "@/providers/user-provider";
import { PreferencesProvider } from "@/providers/preferences-provider";
import { ToastContainer } from "react-toastify";
import { ScrollToTop } from "@/components/shared/scroll-to-top";
import type { ReactNode } from "react";

const jakarta = Plus_Jakarta_Sans({
  variable: "--font-jakarta",
  subsets: ["latin"],
  display: "swap",
});

export const metadata: Metadata = {
  title: "TicketBooking | Find and book live events",
  description:
    "Concerts, festivals and club nights. Book tickets in a few taps.",
};


export default function RootLayout({
  children,
}: {
  children: ReactNode;
}) {
  return (
    <html
      lang="en"
      suppressHydrationWarning
      className={`${jakarta.variable} dark h-full antialiased pretty-scroll`}
    >
      <body
        suppressHydrationWarning
        className="flex min-h-full flex-col bg-bg text-fg pretty-scroll"
      >
        <ThemeProvider>
          <UserProvider>
            <PreferencesProvider>
              <ScrollToTop />
              {children}
            </PreferencesProvider>
          </UserProvider>
        </ThemeProvider>
        <ToastContainer />
      </body>
    </html>
  );
}
