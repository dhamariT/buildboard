'use client';

import Link from 'next/link';
import './earlystart.css';

export default function EarlyStartPage() {
  const handleSlackClick = () => {
    window.open('https://hackclub.slack.com/archives/C09PAAVLZ16e', '_blank', 'noopener,noreferrer');
  };

  return (
    <div className="full">
      <div className="content start" style={{ flexDirection: 'column', gap: '2rem' }}>
        <h1 style={{
          fontFamily: 'Codystar, cursive',
          fontSize: 'clamp(32px, 5vw, 64px)',
          margin: 0
        }}>
          Join the BuildBoard Channel
        </h1>

        <p style={{
          fontSize: 'clamp(18px, 2.5vw, 32px)',
          textAlign: 'center',
          maxWidth: '800px',
          margin: '0 2rem'
        }}>
             for updates and announcements.
        </p>

        <button
          onClick={handleSlackClick}
          className="button"
          style={{
            padding: '1rem 2rem',
            fontSize: 'clamp(18px, 2vw, 28px)',
            fontFamily: 'Codystar, cursive',
            backgroundColor: 'transparent',
            color: 'white',
            border: '2px solid white',
            cursor: 'pointer',
            transition: 'all 0.3s ease'
          }}
          onMouseEnter={(e) => {
            e.target.style.color = 'magenta';
            e.target.style.borderColor = 'magenta';
          }}
          onMouseLeave={(e) => {
            e.target.style.color = 'white';
            e.target.style.borderColor = 'white';
          }}
        >
            Slack Channel
        </button>

        <Link
          href="/"
          style={{
            color: 'white',
            textDecoration: 'underline',
            fontSize: 'clamp(16px, 2vw, 24px)',
            marginTop: '2rem'
          }}
        >
          Back to Home
        </Link>
      </div>
    </div>
  );
}