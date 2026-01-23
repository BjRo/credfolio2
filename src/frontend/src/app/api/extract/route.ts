import { type NextRequest, NextResponse } from "next/server";

const BACKEND_URL = process.env.BACKEND_URL || "http://localhost:8080";

// Increase timeout for long-running extractions
export const maxDuration = 120; // seconds

export async function POST(request: NextRequest) {
  try {
    // Read the entire body as array buffer
    const body = await request.arrayBuffer();
    const contentType = request.headers.get("content-type") || "";

    console.log(`Proxying request: ${body.byteLength} bytes, content-type: ${contentType}`);

    const response = await fetch(`${BACKEND_URL}/api/extract`, {
      method: "POST",
      headers: {
        "Content-Type": contentType,
        "Content-Length": body.byteLength.toString(),
      },
      body: body,
    });

    console.log(`Backend response: ${response.status}`);

    const data = await response.json();

    if (!response.ok) {
      return NextResponse.json(data, { status: response.status });
    }

    return NextResponse.json(data);
  } catch (error) {
    console.error("Proxy error:", error);
    return NextResponse.json({ error: "Failed to connect to backend" }, { status: 502 });
  }
}
