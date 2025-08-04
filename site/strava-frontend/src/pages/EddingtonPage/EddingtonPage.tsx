
import { FC } from "react";
import { Flex, Text } from "@chakra-ui/react";
import { useParams } from "react-router-dom";
import { getAthlete, getRoute } from "../../api/rest";
import { useQuery } from "@tanstack/react-query";
import { NotFound } from "../404/404";
import { Loading } from "../../components/Loading/Loading";
import { useAuthenticated } from "../../contexts/Authenticated";
import { ErrorBox } from "../../components/ErrorBox/ErrorBox";
import { Eddington } from "./Eddington";
import { EddingtonAllChart } from "../../components/Eddington/EddingtonAllChart";

export const EddingtonPage: FC<{}> = ({}) => {
  const { athlete_id } = useParams();
  const { authenticatedUser, isLoading } = useAuthenticated();

  if (isLoading) {
    return <Loading />;
  }

  if (
    authenticatedUser?.athlete_id?.toString() === athlete_id?.toString() ||
    // Or Steven
    authenticatedUser?.athlete_id?.toString() === "2661162"
  ) {
    return <>
      <Eddington />
    </>;
  }
  return <>
    <ErrorBox error="You are not allowed to view this page." />
  </>;
};
