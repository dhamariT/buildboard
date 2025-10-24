'use client';

import { useEffect, useState, useRef } from 'react';
import { FiVolume2, FiVolumeX } from 'react-icons/fi';
import { useRouter } from 'next/navigation';
import { earlyStartAPI } from '@/lib/api';


export default function MusicPlayer() {
  const [isPlaying, setIsPlaying] = useState(false);
  const [verifiedCount, setVerifiedCount] = useState(0);
  const [showWait, setShowWait] = useState(true);
  const [muted, setMuted] = useState(false);

  // Early access signup states
  const [signupMode, setSignupMode] = useState('button'); // 'button', 'email', 'otp'
  const [email, setEmail] = useState('');
  const [otp, setOtp] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const audioRef = useRef(null);
  const router = useRouter();

  useEffect(() => {
    // Set showWait to false after component mounts
    setShowWait(false);

    // Fetch verified user count on mount
    earlyStartAPI.getCount()
      .then(data => setVerifiedCount(data.verified))
      .catch(console.error);

    // Poll for updates every 10 seconds
    const interval = setInterval(() => {
      earlyStartAPI.getCount()
        .then(data => setVerifiedCount(data.verified))
        .catch(console.error);
    }, 10000);

    return () => clearInterval(interval);
  }, []);

  const start = () => {
    setIsPlaying(true);
    if (audioRef.current) {
      audioRef.current.play();
      audioRef.current.loop = true;
    }
  };

  const stop = () => {
    setIsPlaying(false);
    if (audioRef.current) {
      audioRef.current.pause();
      audioRef.current.currentTime = 0;
    }
  };


  const handleEmailSubmit = async (e) => {
    e.preventDefault();
    if (!email) return;

    setLoading(true);
    setError('');

    try {
      await earlyStartAPI.signup({ email });
      setSignupMode('otp');
    } catch (err) {
      setError(err.message || 'Failed to send code');
    } finally {
      setLoading(false);
    }
  };

  const handleOTPSubmit = async (e) => {
    e.preventDefault();
    if (!otp || otp.length !== 6) return;

    setLoading(true);
    setError('');

    try {
      await earlyStartAPI.verifyOTP({ email, otp: otp.toUpperCase() });
      // Redirect to early start page after successful verification
      router.push('/earlystart');
    } catch (err) {
      setError(err.message || 'Invalid code');
    } finally {
      setLoading(false);
    }
  };

  const toggleMute = () => {
    setMuted(prev => {
      const next = !prev;
      if (audioRef.current) {
        audioRef.current.muted = next;
      }
      return next;
    });
  };

  if (showWait) {
    return (
      <>
        <div className="content wait">
          BuildBord: Your Work on a New York Billboard
        </div>
        <button
          className="mute-btn"
          onClick={toggleMute}
          aria-label={muted ? 'Unmute music' : 'Mute music'}
          title={muted ? 'Unmute' : 'Mute'}
        >
          {muted ? <FiVolumeX aria-hidden="true" /> : <FiVolume2 aria-hidden="true" />}
        </button>
      </>
    );
  }

  if (!isPlaying) {
    return (
      <>
        <audio ref={audioRef} src="/texture - 184 (rnb).wav" muted={muted} />
        <div className="content start" onClick={start}>
          Build something you&apos;re proud of, and we&apos;ll put it on a billboard in New York City. Seriously.
        </div>
        <button
          className="mute-btn"
          onClick={toggleMute}
          aria-label={muted ? 'Unmute music' : 'Mute music'}
          title={muted ? 'Unmute' : 'Mute'}
        >
          {muted ? <FiVolumeX aria-hidden="true" /> : <FiVolume2 aria-hidden="true" />}
        </button>
      </>
    );
  }

  return (
    <>
      <audio ref={audioRef} src="/texture - 184 (rnb).wav" muted={muted} />
      <div className="content playing">
        <marquee scrollamount="15">BUILD YOUR PROJECTS</marquee>
        <marquee scrollamount="15" direction="right">50 SPOTS</marquee>

        <marquee
          behavior="alternate"
          direction="right"
          className="count"
        >
          {verifiedCount} PEOPLE STARTING EARLY
        </marquee>

        {signupMode === 'button' && (
          <div
            className="button signup-button"
            onMouseEnter={() => setSignupMode('email')}
          >
            CLAIM YOUR SPOT <span style={{color: 'red'}}>EARLY</span>
          </div>
        )}

        {signupMode === 'email' && (
          <form onSubmit={handleEmailSubmit} className="inline-signup">
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="your@email.com"
              autoFocus
              disabled={loading}
              onBlur={(e) => {
                // Only reset if input is empty and not submitting
                if (!e.target.value && !loading) {
                  setTimeout(() => setSignupMode('button'), 200);
                }
              }}
            />
            {error && <div className="error-text">{error}</div>}
          </form>
        )}

        {signupMode === 'otp' && (
          <form onSubmit={handleOTPSubmit} className="inline-signup">
            <input
              type="text"
              value={otp}
              onChange={(e) => setOtp(e.target.value.toUpperCase())}
              placeholder="ENTER CODE"
              maxLength={6}
              autoFocus
              disabled={loading}
            />
            {error && <div className="error-text">{error}</div>}
          </form>
        )}

        <marquee scrollamount="15">MECHANICAL KEYBOARDS</marquee>
        <marquee scrollamount="15" direction="right">WEB APPS THAT SOLVE PROBLEMS</marquee>

        <div className="button">
          Hardware Projects
        </div>

        <marquee
          behavior="alternate"
          className="tempo"
        >
          PHOTO + QR CODE TO GITHUB
        </marquee>

        <div className="button">
          Software Projects
        </div>

        <marquee scrollamount="15">QR CODE TO LIVE PROJECT</marquee>
        <marquee scrollamount="15" direction="right">YOU BUILT IT NOW LET&apos;S SHOW IT</marquee>
      </div>
      <button
        className="mute-btn"
        onClick={toggleMute}
        aria-label={muted ? 'Unmute music' : 'Mute music'}
        title={muted ? 'Unmute' : 'Mute'}
      >
        {muted ? '🔇' : '🔈'}
      </button>
    </>
  );
}