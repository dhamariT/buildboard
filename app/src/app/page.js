import MusicPlayer from './components/MusicPlayer';

export default function Home() {
  return (
    <>
      <MusicPlayer />
      <footer className="footer">
        Plays back real-time generated 3-part (melody, bass, drums) from a miniaturized{' '}
        <a href="https://g.co/magenta/musicvae" target="_blank" rel="noopener noreferrer">
          MusicVAE
        </a>{' '}
        model. Listen to samples from the full model on{' '}
        <a
          href="https://www.youtube.com/watch?v=xU1W3c9p2RU&list=PLBUMAYA6kvGVds2N7HIMQnZc0SMFk99Yl"
          target="_blank"
          rel="noopener noreferrer"
        >
          YouTube
        </a>{' '}
        or sample your own in{' '}
        <a
          href="https://goo.gl/magenta/musicvae-colab"
          target="_blank"
          rel="noopener noreferrer"
        >
          Colab
        </a>
        .
      </footer>
    </>
  );
}