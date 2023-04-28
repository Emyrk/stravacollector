import axios, { AxiosError, AxiosPromise, AxiosResponse } from "axios"
import * as TypesGen from "./typesGenerated"

export type ApiError = AxiosError<TypesGen.Response> & {
  response: AxiosResponse<TypesGen.Response>
}

export const getAuthenticatedUser = async (): Promise<
  TypesGen.AthleteSummary | undefined
> => {
  try {
    const response = await axios.get<TypesGen.AthleteSummary>("/api/v1/whoami")
    return response.data
  } catch (error) {
    throw error
  }
}

export const isAPIError = (err: unknown): err is ApiError => {
  if (axios.isAxiosError(err)) {
    const response = err.response?.data
    if (!response) {
      return false
    }

    return (
      typeof response.message === "string" &&
      (typeof response.errors === "undefined" || Array.isArray(response.errors))
    )
  }

  return false
}

export const toAPIError = (err: unknown): ApiError => {
  if(!isAPIError(err) || !err) {
    throw new Error("not an API error")
  }
  const ax = err as AxiosError<ApiError>
  if(!ax.response) {
    throw new Error("no response data")
  }
  return ax.response.data
}

/**
 *
 * @param error
 * @param defaultMessage
 * @returns error's message if ApiError or Error, else defaultMessage
 */
export const getErrorMessage = (
  error: Error | ApiError | unknown,
  defaultMessage: string,
): string =>
  isAPIError(error)
    ? error.response.data.message
    : error instanceof Error
    ? error.message
    : defaultMessage

export const getErrorDetail = (
  error: Error | ApiError | unknown,
): string | undefined | null =>
  isAPIError(error)
    ? error.response.data.detail
    : error instanceof Error
    ? error.stack
    : null