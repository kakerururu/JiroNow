import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './index.css';
import App from './App.tsx';

// Setup MSW in development mode
if (process.env.NODE_ENV === 'development') {
  const { setupWorker } = await import('msw/browser');
  const { getJiroNowAPIMock } = await import('./api/jiroNowAPI.msw.ts');

  const worker = setupWorker(...getJiroNowAPIMock());
  worker.start();
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
