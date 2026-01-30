import { type NextRequest, NextResponse } from "next/server";

const MINIO_INTERNAL_URL = process.env.MINIO_INTERNAL_URL || "http://credfolio2-minio:9000";

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params;
  const pathString = path.join("/");

  // Get query parameters (for presigned URL params)
  const searchParams = request.nextUrl.searchParams.toString();
  const queryString = searchParams ? `?${searchParams}` : "";

  // Construct the internal MinIO URL
  const minioUrl = `${MINIO_INTERNAL_URL}/credfolio/${pathString}${queryString}`;

  try {
    const rangeHeader = request.headers.get("range");
    const response = await fetch(minioUrl, {
      headers: {
        // Forward range requests for partial content
        ...(rangeHeader ? { range: rangeHeader } : {}),
      },
    });

    if (!response.ok) {
      return new NextResponse(null, { status: response.status });
    }

    // Get the content type and body
    const contentType = response.headers.get("content-type") || "application/octet-stream";
    const contentLength = response.headers.get("content-length");
    const body = await response.arrayBuffer();

    return new NextResponse(body, {
      status: response.status,
      headers: {
        "Content-Type": contentType,
        "Cache-Control": "public, max-age=31536000, immutable",
        ...(contentLength ? { "Content-Length": contentLength } : {}),
      },
    });
  } catch (error) {
    console.error("Storage proxy error:", error);
    return new NextResponse(null, { status: 502 });
  }
}
