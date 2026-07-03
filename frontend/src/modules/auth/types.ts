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

export type LoginResponse = AuthUser

export interface RegisterRequest {
  email: string
  password: string
  name: string
}
