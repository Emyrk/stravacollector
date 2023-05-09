import axios, { AxiosError, AxiosPromise, AxiosResponse } from "axios"
import * as TypesGen from "./typesGenerated"
// import JSONbig from 'json-bigint';

// Overriding the transformResponse of axios and converting any number which crosses JS max limit to string using stringify
// Remember the data is received as string in case of string not JSON over the network that's why we need parser always
// Default JSON.parse will transform the huge number to some random number which is an issue
// axios.defaults.transformResponse = [(data) => {
//   try {
//     return JSONbig.parse(data);
//   } catch (error) {
//     return data;
//   }
// }];


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

export const getHugelLeaderBoard = async (): Promise<
  TypesGen.HugelLeaderBoard | undefined
> => {
  try {
    const response = await axios.get<TypesGen.HugelLeaderBoard>("/api/v1/hugelboard")
    return response.data
  } catch (error) {
    throw error
  }
}

export const getSuperHugelLeaderBoard = async (): Promise<
  TypesGen.SuperHugelLeaderBoard | undefined
> => {
  try {
    const response = await axios.get<TypesGen.SuperHugelLeaderBoard>("/api/v1/superhugelboard")
    return response.data
  } catch (error) {
    throw error
  }
}

export const getRoute = async (
  routeName: string
): Promise<
  TypesGen.CompetitiveRoute | undefined
> => {
  try {
    const response = await axios.get<TypesGen.CompetitiveRoute>(`/api/v1/route/${routeName}`)
    return response.data
  } catch (error) {
    throw error
  }
}

export const getDetailedSegments = async (
  segments: string[]
): Promise<
  TypesGen.PersonalSegment[] | undefined
> => {
  try {
    const response = await axios.post<TypesGen.PersonalSegment[]>(`/api/v1/segments`, segments)
    return response.data
  } catch (error) {
    throw error
  }
}

export const getHugelSegments = async (): Promise<
  TypesGen.CompetitiveRoute | undefined
> => {
  try {
    const response = await axios.get<TypesGen.CompetitiveRoute>("/api/v1/route/das-hugel")
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