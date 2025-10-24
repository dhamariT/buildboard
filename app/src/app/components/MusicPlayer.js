'use client';

import { useEffect, useState, useRef } from 'react';
import { FiVolume2, FiVolumeX } from 'react-icons/fi';
import { useRouter } from 'next/navigation';


export default function MusicPlayer() {
  const [isPlaying, setIsPlaying] = useState(false);
  const [count, setCount] = useState(0);
  const [showWait, setShowWait] = useState(true);
  const [muted, setMuted] = useState(false);

  const audioRef = useRef(null);
  const router = useRouter();

  useEffect(() => {
    // Set showWait to false after component mounts
    setShowWait(false);
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

  const nextProject = () => {
    setCount(prev => prev + 1);
    if (audioRef.current) {
      audioRef.current.currentTime = 0;
      audioRef.current.play();
    }
  };

  const claimEarly = () => {
    // Stop audio and navigate to the guide page
    stop();
    router.push('/earlystart/guide');
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
        <marquee scrollAmount={15}>BUILD YOUR PROJECTS</marquee>
        <marquee scrollAmount={15} direction="right">50 SPOTS FOR TEENAGERS</marquee>

        <div className="button" onClick={nextProject}>
          Next Project
        </div>

        <marquee
          behavior="alternate"
          direction="right"
          className="count"
        >
          {count} PROJECTS SHIPPED
        </marquee>

          <div className="button" onClick={claimEarly}>
              CLAIM YOUR SPOT <span style={{color: 'red'}}>EARLY</span>
          </div>

        <marquee scrollAmount={15}>MECHANICAL KEYBOARDS</marquee>
        <marquee scrollAmount={15} direction="right">WEB APPS THAT SOLVE PROBLEMS</marquee>

        <div className="button">
          Physical Projects
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

        <marquee scrollAmount={15}>QR CODE TO LIVE PROJECT</marquee>
        <marquee scrollAmount={15} direction="right">YOU BUILT IT NOW LET&apos;S SHOW IT</marquee>
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