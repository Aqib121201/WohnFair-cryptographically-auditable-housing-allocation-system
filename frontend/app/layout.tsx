import './globals.css';
import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'WohnFair',
  description: 'Fairness-preserving, verifiable housing allocation',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
