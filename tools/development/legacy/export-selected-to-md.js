const fs = require("fs");
const path = require("path");

const args = process.argv.slice(2);

let outputFile = "SELECTED_CODE_CONTEXT.md";
const targets = [];

for (let i = 0; i < args.length; i++) {
  if (args[i] === "-o" || args[i] === "--output") {
    outputFile = args[i + 1];
    i++;
  } else {
    targets.push(args[i]);
  }
}

if (targets.length === 0) {
  console.error("Cách dùng:");
  console.error("node export-selected-to-md.js -o OUTPUT.md folder1 folder2 folder3");
  process.exit(1);
}

const ROOT = process.cwd();
const outputPath = path.resolve(outputFile);

const IGNORE_DIRS = new Set([
  "node_modules",
  ".git",
  ".next",
  "dist",
  "build",
  "coverage",
  ".vercel",
  ".turbo",
  ".cache",
  ".idea",
  ".vscode",
  "tmp",
  "temp",
  "logs",
  "bin",
  "vendor"
]);

const IGNORE_FILES = new Set([
  ".DS_Store",
  "package-lock.json",
  "yarn.lock",
  "pnpm-lock.yaml",
  "bun.lockb"
]);

const IGNORE_EXTENSIONS = new Set([
  ".png",
  ".jpg",
  ".jpeg",
  ".gif",
  ".webp",
  ".ico",
  ".svg",
  ".mp4",
  ".mp3",
  ".wav",
  ".woff",
  ".woff2",
  ".ttf",
  ".eot",
  ".zip",
  ".rar",
  ".7z",
  ".tar",
  ".gz",
  ".pdf",
  ".exe",
  ".dll",
  ".so",
  ".dylib",
  ".db",
  ".sqlite",
  ".sqlite3"
]);

const MAX_FILE_SIZE_KB = 800;

function normalizePath(filePath) {
  return filePath.replace(/\\/g, "/");
}

function isEnvSecretFile(name) {
  return (
    name === ".env" ||
    name.startsWith(".env.") ||
    name.endsWith(".env")
  );
}

function shouldIgnore(fullPath, name) {
  const stat = fs.statSync(fullPath);

  if (stat.isDirectory()) {
    return IGNORE_DIRS.has(name);
  }

  if (IGNORE_FILES.has(name)) return true;

  // Không đưa file env thật vào prompt để tránh lộ secret.
  // Các .env develop vẫn chứa key runtime nên không xuất vào tài liệu.
  if (isEnvSecretFile(name)) return true;

  const ext = path.extname(name).toLowerCase();
  if (IGNORE_EXTENSIONS.has(ext)) return true;

  const sizeKB = stat.size / 1024;
  if (sizeKB > MAX_FILE_SIZE_KB) return true;

  return false;
}

function getLanguage(fileName) {
  const lower = fileName.toLowerCase();
  const ext = path.extname(lower);

  if (lower === "dockerfile") return "dockerfile";
  if (lower === "makefile") return "makefile";

  const map = {
    ".go": "go",
    ".mod": "go",
    ".sum": "text",
    ".proto": "protobuf",
    ".sql": "sql",
    ".sh": "bash",
    ".bash": "bash",
    ".js": "javascript",
    ".jsx": "jsx",
    ".ts": "typescript",
    ".tsx": "tsx",
    ".json": "json",
    ".yaml": "yaml",
    ".yml": "yaml",
    ".toml": "toml",
    ".xml": "xml",
    ".html": "html",
    ".css": "css",
    ".scss": "scss",
    ".md": "markdown",
    ".txt": "text"
  };

  return map[ext] || "text";
}

function sortItems(dir, items) {
  return items.sort((a, b) => {
    const aPath = path.join(dir, a);
    const bPath = path.join(dir, b);

    const aIsDir = fs.statSync(aPath).isDirectory();
    const bIsDir = fs.statSync(bPath).isDirectory();

    if (aIsDir && !bIsDir) return -1;
    if (!aIsDir && bIsDir) return 1;

    return a.localeCompare(b);
  });
}

function getVisibleItems(dir) {
  return sortItems(
    dir,
    fs.readdirSync(dir).filter((item) => {
      const fullPath = path.join(dir, item);
      return !shouldIgnore(fullPath, item);
    })
  );
}

function buildTree(dir, prefix = "") {
  let result = "";
  const items = getVisibleItems(dir);

  items.forEach((item, index) => {
    const fullPath = path.join(dir, item);
    const isLast = index === items.length - 1;

    result += prefix + (isLast ? "└── " : "├── ") + item + "\n";

    if (fs.statSync(fullPath).isDirectory()) {
      result += buildTree(fullPath, prefix + (isLast ? "    " : "│   "));
    }
  });

  return result;
}

function collectFiles(dir, files = []) {
  const items = getVisibleItems(dir);

  for (const item of items) {
    const fullPath = path.join(dir, item);
    const stat = fs.statSync(fullPath);

    if (stat.isDirectory()) {
      collectFiles(fullPath, files);
    } else {
      files.push(fullPath);
    }
  }

  return files;
}

function readFileSafe(filePath) {
  try {
    return fs.readFileSync(filePath, "utf8");
  } catch (error) {
    return `Không thể đọc file: ${error.message}`;
  }
}

function generateFolderSection(targetInput) {
  const targetPath = path.resolve(targetInput);

  if (!fs.existsSync(targetPath)) {
    return `## Không tìm thấy thư mục: ${targetInput}\n\n`;
  }

  const rootName = path.basename(targetPath);
  const files = collectFiles(targetPath);

  let md = "";

  md += `---\n\n`;
  md += `# Folder: ${normalizePath(path.relative(ROOT, targetPath))}\n\n`;

  md += `## Folder Structure\n\n`;
  md += "```txt\n";
  md += `${rootName}\n`;
  md += buildTree(targetPath);
  md += "```\n\n";

  md += `## Files And Code\n\n`;

  for (const filePath of files) {
    const relativePath = normalizePath(path.relative(ROOT, filePath));
    const fileName = path.basename(filePath);
    const lang = getLanguage(fileName);
    const content = readFileSafe(filePath);

    md += `### ${relativePath}\n\n`;
    md += "```" + lang + "\n";
    md += content;
    if (!content.endsWith("\n")) md += "\n";
    md += "```\n\n";
  }

  return md;
}

function generateMarkdown() {
  let md = "";

  md += `# Selected Project Code Context\n\n`;
  md += `Generated at: ${new Date().toISOString()}\n\n`;

  md += `## Project Root\n\n`;
  md += "```txt\n";
  md += `${normalizePath(ROOT)}\n`;
  md += "```\n\n";

  md += `## Exported Folders\n\n`;
  md += "```txt\n";
  for (const target of targets) {
    md += `${normalizePath(target)}\n`;
  }
  md += "```\n\n";

  md += `## Ignored\n\n`;
  md += "```txt\n";
  md += "Dirs:\n";
  md += Array.from(IGNORE_DIRS).join("\n");
  md += "\n\nFiles:\n";
  md += Array.from(IGNORE_FILES).join("\n");
  md += "\n\nExtensions:\n";
  md += Array.from(IGNORE_EXTENSIONS).join("\n");
  md += "\n\nEnv secrets:\n.env, .env.*, *.env";
  md += "\n```\n\n";

  for (const target of targets) {
    md += generateFolderSection(target);
  }

  return md;
}

try {
  const markdown = generateMarkdown();
  fs.writeFileSync(outputPath, markdown, "utf8");

  console.log("DONE: Đã xuất code ra file:");
  console.log(outputPath);
} catch (error) {
  console.error("ERROR:", error.message);
  process.exit(1);
}
