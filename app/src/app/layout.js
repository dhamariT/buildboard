import "./globals.css";

export const metadata = {
  title: "BuildBord - Your Work on a New York Billboard",
  description: "Build something you're proud of, and we'll put it on a billboard in New York City. 50 spots for teenagers who ship projects they actually care about.",
  keywords: ["teenage builders", "NYC billboard", "student projects", "showcase work", "coding projects"],
  openGraph: {
    title: "BuildBord - Your Work on a New York Billboard",
    description: "Build something you're proud of, and we'll put it on a billboard in New York City.",
    type: "website",
  },
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <head>
        <link href='https://fonts.googleapis.com/css?family=Codystar' rel='stylesheet' />
      </head>
      <body>
        {children}
      </body>
    </html>
  );
}