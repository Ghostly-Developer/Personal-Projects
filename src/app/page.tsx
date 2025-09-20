'use client';
import React, { useState } from "react";
import Header from "@/components/header";
import Footer from "@/components/footer";

export default function Home() {
  const [board, setBoard] = useState(Array(9).fill(""));

  const handleClick = (index: number) => {
    if (board[index]) return;
    const newBoard = [...board];
    newBoard[index] = "X"; // For now, always "X"
    setBoard(newBoard);
  };

  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-grow flex items-center justify-center">
        <div className="container mx-auto px-4">
          <div className="text-center">
            <h2 className="text-3xl font-bold mb-8">Welcome to Tic Tac Toe</h2>
            <div className="tic-board">
              {board.map((cell, idx) => (
                <button
                  key={idx}
                  className="tic-cell"
                  onClick={() => handleClick(idx)}
                >
                  {cell}
                </button>
              ))}
            </div>
          </div>
        </div>
      </main>
      <Footer />
    </div>
  );
}
