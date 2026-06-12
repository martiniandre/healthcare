import { Component } from "react"
import type { ErrorInfo, ReactNode } from "react"
import i18next from "i18next"
import { AlertTriangle } from "lucide-react"
import { Card } from "./ui/Card"
import { Button } from "./ui/Button"

interface Props {
  children: ReactNode
}

interface State {
  hasError: boolean
  error: Error | null
}

export class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    hasError: false,
    error: null,
  }

  public static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error }
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error("Uncaught error in React component:", error, errorInfo)
  }

  private handleReset = () => {
    this.setState({ hasError: false, error: null })
  }

  public render() {
    if (this.state.hasError) {
      return (
        <div className="flex w-full items-center justify-center p-6 bg-background">
          <Card glowingType="cyan" className="p-6 max-w-lg w-full flex flex-col items-center text-center gap-4 border-red-100">
            <div className="bg-red-50 p-3 rounded-full mb-2">
              <AlertTriangle className="w-8 h-8 text-red-500" />
            </div>
            <h2 className="text-lg font-bold text-gray-900">
              {i18next.t("header:errorBoundary.title")}
            </h2>
            <p className="text-sm text-gray-500 max-w-sm">
              {i18next.t("header:errorBoundary.description")}
            </p>
            <div className="bg-gray-50 rounded-lg p-3 w-full text-left overflow-hidden border border-gray-100">
              <code className="text-xs text-red-600 block break-words whitespace-pre-wrap font-mono">
                {this.state.error?.message || "Unknown Error"}
              </code>
            </div>
            <Button onClick={this.handleReset} className="mt-2 bg-gray-200 text-gray-800 hover:bg-gray-300 border-none">
              {i18next.t("header:errorBoundary.retry")}
            </Button>
          </Card>
        </div>
      )
    }

    return this.props.children
  }
}
