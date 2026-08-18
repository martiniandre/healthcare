export interface LoginRequest {
  email: string
  password: string
}

export interface AuthUser {
  userId: string
  role: string
  email: string
  fullName?: string
  isActive?: boolean
}

export interface LoginResponse extends AuthUser {
  csrfToken?: string
}

export interface RegisterRequest {
  email: string
  password: string
  name: string
}
