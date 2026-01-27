import type { Metadata } from "next";
import { SiteHeader } from "@/components/site-header";
import { ThemeProvider } from "@/components/theme-provider";
import { UrqlProvider } from "@/lib/urql";
import "./globals.css";

export const metadata: Metadata = {
  title: "Credfolio",
  description: "Your professional portfolio powered by AI",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className="antialiased">
        <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
          <SiteHeader />
          <UrqlProvider>{children}</UrqlProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
