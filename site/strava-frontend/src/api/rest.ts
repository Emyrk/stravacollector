import axios from "axios"
import * as TypesGen from "./typesGenerated"

export const getAuthenticatedUser = async (): Promise<
  TypesGen.AthleteLogin | undefined
> => {
  try {
    const response = await axios.get<TypesGen.AthleteLogin>("/api/v1/whoami")
    return response.data
  } catch (error) {
    if (axios.isAxiosError(error) && error.response?.status === 401) {
      return undefined
    }

    throw error
  }
}