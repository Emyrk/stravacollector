import { FC, useState } from "react";
import { useParams } from "react-router-dom";
import { getAthleteEddington } from "../../api/rest";
import { Eddington } from "../../api/typesGenerated";
import { useQuery } from "@tanstack/react-query";
import { Loading } from "../../components/Loading/Loading";
import { ErrorBox } from "../../components/ErrorBox/ErrorBox";




export const EddingtonChart: FC<{}> = ({}) => {
  const { athlete_id } = useParams();
  const [ eddington, setEddington] = useState<Eddington>();


  const queryKey = ["athlete", athlete_id, "eddington"];
    const {
      data: chartData,
      error: chartError,
      isLoading: chartLoading,
      isFetched: chartFetched,
    } = useQuery({
      queryKey,
      enabled: !!athlete_id,
      queryFn: () =>
        getAthleteEddington(athlete_id || "me"),
      onSuccess: (data) => {
        setEddington(data)
      },
      onError: (error) => {
        console.error("Error fetching athlete data:", error);
      }
    });

  if (
    (!chartData || chartLoading)
  ) {
    return <Loading />;
  }

  if (
    chartError
  ) {
    return <ErrorBox error="Error fetching eddington data." detail={chartError} />;
  }
    

  return (
    <>

    </>
  );
};

