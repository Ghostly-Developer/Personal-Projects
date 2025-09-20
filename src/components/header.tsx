import Link from 'next/link';
import Image from 'next/image';

const Header = () => {
  return (
    <header className="w-full bg-zinc-900 text-white shadow-lg">
      <div className="container mx-auto px-4 py-4">
        <nav className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <Link href="/" className="flex items-center space-x-2">
              <Image src="/logo.png" alt="Logo" width={40} height={40} className="rounded-lg" />
              <span className="text-xl font-bold">Ghostly Games</span>
            </Link>
          </div>
          
          <div className="flex items-center space-x-6">
            <Link 
              href="/games/tictactoe"
              className="font-semibold hover:text-zinc-300 transition-colors"
            >
              Tic Tac Toe
            </Link>
            <Link 
              href="/games/mathic"
              className="font-semibold hover:text-zinc-300 transition-colors"
            >
              Mathic
            </Link>
            <Link 
              href="/room"
              className="px-4 py-2 bg-slate-100 hover:bg-slate-300 text-black font-semibold rounded-lg transition-colors shadow-sm"
            >
              Join Room
            </Link>
          </div>
        </nav>
      </div>
    </header>
  );
};

export default Header;
