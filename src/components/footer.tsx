import GhostlyLogo from '@/assets/GhostlyLogo.png';
import Image from 'next/image';

const Footer = () => {
  return (
    <footer className="w-full bg-zinc-900 text-white mt-auto">
      <div className="container mx-auto px-4 py-6">
        <div className="relative">
          <div className="flex flex-col items-center text-center">
            <h3 className="text-lg font-bold mb-2">Ghostly Games</h3>
            <p className="text-zinc-400">Challenge your friends in exciting multiplayer games!</p>
            <div className="flex items-center justify-center gap-2 text-zinc-500 text-sm mt-4">
              © {new Date().getFullYear()} <Image src={GhostlyLogo} alt="Ghostly Logo" width={20} height={20} /> Ghostly Developer. All rights reserved.
            </div>
          </div>
          
          <div className="absolute top-0 right-0 flex items-center space-x-8">
            <a 
              href="https://github.com/Ghostly-Developer" 
              target="_blank" 
              rel="noopener noreferrer"
              className="text-zinc-400 hover:text-white transition-colors hover:underline"
            >
              GitHub
            </a>
            <a 
              href="/about_me" 
              className="text-zinc-400 hover:text-white transition-colors hover:underline"
            >
              About Me
            </a>
            <a 
              href="/contact_us" 
              className="text-zinc-400 hover:text-white transition-colors hover:underline"
            >
              Contact Us
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
