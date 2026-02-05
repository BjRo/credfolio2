import Link from "next/link";
import { ThemeToggle } from "./theme-toggle";

export function SiteHeader() {
  return (
    <header className="bg-foreground text-background">
      <div className="mx-auto flex h-14 max-w-5xl items-center justify-between px-4">
        <Link href="/" className="text-lg font-semibold">
          Credfolio
        </Link>
        <nav className="flex items-center gap-4">
          <Link
            href="/upload"
            className="text-sm text-background/80 hover:text-background transition-colors"
          >
            Upload
          </Link>
          <ThemeToggle invertColors />
        </nav>
      </div>
    </header>
  );
}
