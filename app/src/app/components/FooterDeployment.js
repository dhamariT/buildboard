"use client";

import { useEffect, useState } from "react";

const REPO_OWNER = "dhamariT";
const REPO_NAME = "buildboard";

function formatTimeAgo(iso) {
  if (!iso) return "unknown time";
  const date = new Date(iso);
  const diff = Date.now() - date.getTime();
  const seconds = Math.floor(diff / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);
  if (days > 0) return `${days} day${days === 1 ? "" : "s"} ago`;
  if (hours > 0) return `${hours} hour${hours === 1 ? "" : "s"} ago`;
  if (minutes > 0) return `${minutes} minute${minutes === 1 ? "" : "s"} ago`;
  return `${seconds} second${seconds === 1 ? "" : "s"} ago`;
}

export default function FooterDeployment() {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [info, setInfo] = useState({
    commitSha: "",
    commitUrl: "",
    deployedAt: "",
    runUrl: "",
    source: "",
  });

  useEffect(() => {
    let cancelled = false;

    async function fetchInfo() {
      setLoading(true);
      setError("");
      try {
        // Try latest successful workflow run (proxy for deployment)
        const runsResp = await fetch(
          `https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/actions/runs?status=success&per_page=1`
        );
        if (runsResp.ok) {
          const runs = await runsResp.json();
          const run = runs?.workflow_runs?.[0];
          if (run) {
            const commitSha = run.head_sha || "";
            if (!cancelled) {
              setInfo({
                commitSha,
                commitUrl: commitSha
                  ? `https://github.com/${REPO_OWNER}/${REPO_NAME}/commit/${commitSha}`
                  : `https://github.com/${REPO_OWNER}/${REPO_NAME}`,
                deployedAt: run.updated_at || run.created_at || "",
                runUrl: run.html_url,
                source: "actions",
              });
              setLoading(false);
              return;
            }
          }
        }

        // Fallback: latest commit on the default branch
        const repoResp = await fetch(
          `https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}`
        );
        let defaultBranch = "main";
        if (repoResp.ok) {
          const repo = await repoResp.json();
          defaultBranch = repo.default_branch || defaultBranch;
        }
        const commitsResp = await fetch(
          `https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/commits?sha=${defaultBranch}&per_page=1`
        );
        if (!commitsResp.ok) throw new Error("Failed to load commits");
        const commits = await commitsResp.json();
        const commit = commits?.[0];
        const commitSha = commit?.sha || "";
        if (!cancelled) {
          setInfo({
            commitSha,
            commitUrl: commitSha
              ? `https://github.com/${REPO_OWNER}/${REPO_NAME}/commit/${commitSha}`
              : `https://github.com/${REPO_OWNER}/${REPO_NAME}`,
            deployedAt: commit?.commit?.committer?.date || commit?.commit?.author?.date || "",
            runUrl: `https://github.com/${REPO_OWNER}/${REPO_NAME}`,
            source: "commit",
          });
          setLoading(false);
        }
      } catch (e) {
        if (!cancelled) {
          setError("Could not load deployment info");
          setLoading(false);
        }
      }
    }

    fetchInfo();
    const id = setInterval(fetchInfo, 60_000 * 5); // refresh every 5 minutes
    return () => {
      cancelled = true;
      clearInterval(id);
    };
  }, []);

  const shortSha = info.commitSha ? info.commitSha.slice(0, 7) : "";

  return (
    <span>
      {loading && !error && "Loading deployment…"}
      {!loading && error && (
        <>
          Deployment: unknown · <a href={`https://github.com/${REPO_OWNER}/${REPO_NAME}`} target="_blank" rel="noopener noreferrer">view repo</a>
        </>
      )}
      {!loading && !error && (
        <>
          Last deployment: {info.deployedAt ? `${formatTimeAgo(info.deployedAt)}` : "unknown"}
          {" "}· Commit {shortSha ? (
            <a href={info.commitUrl} target="_blank" rel="noopener noreferrer">{shortSha}</a>
          ) : (
            "unknown"
          )}
          {info.source === "actions" && info.runUrl ? (
            <>
              {" "}· <a href={info.runUrl} target="_blank" rel="noopener noreferrer">workflow</a>
            </>
          ) : null}
        </>
      )}
    </span>
  );
}
