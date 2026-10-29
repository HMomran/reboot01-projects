import { createReadStream, stat } from "node:fs";
import { createServer } from "node:http";
import { extname, resolve, sep } from "node:path";

const root = resolve(".");
const port = Number(process.env.PORT) || 3000;
const types = {
  ".css": "text/css; charset=utf-8",
  ".html": "text/html; charset=utf-8",
  ".ico": "image/x-icon",
  ".js": "text/javascript; charset=utf-8",
  ".json": "application/json; charset=utf-8",
  ".svg": "image/svg+xml",
};

createServer((request, response) => {
  let pathname;
  try {
    pathname = decodeURIComponent(new URL(request.url, "http://localhost").pathname);
  } catch {
    response.writeHead(400).end("Bad request");
    return;
  }
  if (pathname.split("/").some((part) => part.startsWith("."))) {
    response.writeHead(403).end("Forbidden");
    return;
  }
  const file = resolve(root, `.${pathname === "/" ? "/index.html" : pathname}`);

  if (file !== root && !file.startsWith(`${root}${sep}`)) {
    response.writeHead(403).end("Forbidden");
    return;
  }

  stat(file, (error, info) => {
    if (error || !info.isFile()) {
      response.writeHead(404).end("Not found");
      return;
    }

    response.writeHead(200, { "Content-Type": types[extname(file)] || "application/octet-stream" });
    if (request.method === "HEAD") response.end();
    else createReadStream(file).pipe(response);
  });
}).listen(port, () => {
  console.log(`Mini Framework running at http://localhost:${port}`);
});
