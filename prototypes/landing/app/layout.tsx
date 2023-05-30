import { cn } from '@/lib/utils';
import { Inter as FontSans } from 'next/font/google';
import localFont from 'next/font/local';

import 'styles/globals.css';

const fontSans = FontSans({
  subsets: ['latin'],
  variable: '--font-sans',
});

const fontSansBold = localFont({
  weight: '700',
  variable: '--font-sans-bold',
  src: '../assets/fonts/Inter-Bold.ttf',
});

// Font files can be colocated inside of `pages`
const fontHeading = localFont({
  src: '../assets/fonts/CalSans-SemiBold.woff2',
  variable: '--font-heading',
});

export const metadata = {
  title: 'Landing page',
  description: 'Landing page prototype',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body
        className={cn(
          'min-h-screen bg-background font-sans antialiased',
          fontSans.variable,
          fontHeading.variable,
          fontSansBold.variable
        )}
      >
        {children}
      </body>
    </html>
  );
}
