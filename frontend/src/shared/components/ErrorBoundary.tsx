import { Component } from "react"
import type { ErrorInfo, ReactNode } from "react"
import i18next from "i18next"
import { AlertTriangle } from "lucide-react"
import { Card, CardContent } from "./ui/card"
import { Button } from "./ui/button"

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
          <Card className="w-full max-w-lg border-destructive/30">
            <CardContent className="flex flex-col items-center gap-4 text-center">
              <div className="mb-2 rounded-full bg-destructive/10 p-3">
                <AlertTriangle className="h-8 w-8 text-destructive" />
              </div>
              <h2 className="text-lg font-bold text-foreground">
                {i18next.t("header:errorBoundary.title")}
              </h2>
              <p className="max-w-sm text-sm text-muted-foreground">
                {i18next.t("header:errorBoundary.description")}
              </p>
              <div className="w-full overflow-hidden rounded-lg border bg-muted p-3 text-left">
                <code className="block whitespace-pre-wrap break-words font-mono text-xs text-destructive">
                  {this.state.error?.message || "Unknown Error"}
                </code>
              </div>
              <Button variant="secondary" onClick={this.handleReset} className="mt-2">
                {i18next.t("header:errorBoundary.retry")}
              </Button>
            </CardContent>
          </Card>
        </div>
      )
    }

    return this.props.children
  }
}
