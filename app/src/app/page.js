import MusicPlayer from './components/MusicPlayer';
import FooterDeployment from './components/FooterDeployment';

export default function Home() {
  return (
    <>
      <MusicPlayer />
      <footer className="footer">
        <FooterDeployment />
      </footer>
    </>
  );
}