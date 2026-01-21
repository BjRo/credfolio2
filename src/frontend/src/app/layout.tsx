import type { Metadata } from "next";
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
    <html lang="en">
      <body className="antialiased">
        <UrqlProvider>{children}</UrqlProvider>
      </body>
    </html>
  );
}
