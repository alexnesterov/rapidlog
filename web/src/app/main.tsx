import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { QueryClientProvider } from '@tanstack/react-query'
import { RouterProvider } from '@tanstack/react-router'
import './styles/index.css'
import { sessionQuery } from '@/entities/session'
import { createQueryClient } from './providers/queryClient'
import { createAppRouter } from './router'

const queryClient = createQueryClient(() => {
  // сессии больше нет: guard страницы входа пропустит, журнал не вернёт обратно
  queryClient.setQueryData(sessionQuery.queryKey, null)
  void router.navigate({ to: '/login', search: {} })
})
const router = createAppRouter(queryClient)

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </StrictMode>,
)
