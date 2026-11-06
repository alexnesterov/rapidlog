import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { QueryClientProvider } from '@tanstack/react-query'
import './styles/index.css'
import { JournalPage } from '@/pages/journal'
import { queryClient } from './providers/queryClient'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <JournalPage />
    </QueryClientProvider>
  </StrictMode>,
)
