import { FC, PropsWithChildren, createContext, useContext } from "react";
import * as TypesGen from "../api/typesGenerated";
import { getAuthenticatedUser } from "../api/rest";
import { useQuery } from "@tanstack/react-query";

interface AuthenticatedContextValue {
  authenticatedUser?: TypesGen.AthleteLogin
  fetchError?: Error | unknown
  isFetched: boolean
  isLoading: boolean
}

// Undefined is used to indicate that the context is not yet initialized
export const AuthenticatedContext = createContext<AuthenticatedContextValue>({
  isFetched: false,
  isLoading: false,
})

export const AuthenticatedProvider: FC<PropsWithChildren> = ({ children }) => {
  const queryKey = ["authenticated-user"]
  const {
    data: userLogin,
    error: userLoginError,
    isLoading: userLoading,
    isFetched: userFetched,
  } = useQuery({
    queryKey: queryKey,
    queryFn: getAuthenticatedUser,
  })


  return <AuthenticatedContext.Provider value={{
    authenticatedUser: userLogin,
    isFetched: userFetched,
    isLoading: userLoading,
    fetchError: userLoginError,
  }}>
    {children}
  </AuthenticatedContext.Provider>
}

export const useAuthenticated = (): AuthenticatedContextValue => {
  const context = useContext(AuthenticatedContext)

  if (!context) {
    throw new Error("useAuthenticated should be used inside of <AuthenticatedProvider />")
  }

  return context
}