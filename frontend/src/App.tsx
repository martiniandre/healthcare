import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { AppRoutes } from "./app/routes"
import { ToastContainer } from "./shared/components/ui/ToastContainer"
import { useAuthInit } from "./shared/hooks/useAuthInit"
import { Spinner } from "./shared/components/ui/Spinner"

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
})

function AppBootstrap() {
  const { isLoading } = useAuthInit()

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Spinner className="w-8 h-8 text-primary" />
      </div>
    )
  }

  return (
    <>
      <AppRoutes />
      <ToastContainer />
    </>
  )
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AppBootstrap />
    </QueryClientProvider>
  )
}

export default App

