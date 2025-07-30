import { FC, useState } from "react";
import { useParams } from "react-router-dom";
import { getAthleteEddington } from "../../api/rest";
import { Eddington } from "../../api/typesGenerated";
import { useQuery } from "@tanstack/react-query";
import { Loading } from "../../components/Loading/Loading";
import { ErrorBox } from "../../components/ErrorBox/ErrorBox";
import React, { PureComponent } from 'react';
import { BarChart, Bar, Rectangle, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { ContentType } from "recharts/types/component/Tooltip";
import { NameType, ValueType } from "recharts/types/component/DefaultTooltipContent";


const data = [
  {
    name: 'Page A',
    uv: 4000,
    pv: 2400,
    amt: 2400,
  },
  {
    name: 'Page B',
    uv: 3000,
    pv: 1398,
    amt: 2210,
  },
  {
    name: 'Page C',
    uv: 2000,
    pv: 9800,
    amt: 2290,
  },
  {
    name: 'Page D',
    uv: 2780,
    pv: 3908,
    amt: 2000,
  },
  {
    name: 'Page E',
    uv: 1890,
    pv: 4800,
    amt: 2181,
  },
  {
    name: 'Page F',
    uv: 2390,
    pv: 3800,
    amt: 2500,
  },
  {
    name: 'Page G',
    uv: 3490,
    pv: 4300,
    amt: 2100,
  },
];

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
    
  console.log("Eddington Chart Data", chartData);

  return (
    <>
        <BarChart
          width={500}
          height={300}
          data={chartData.miles_histogram.map((value, index) => ({
            index, value
          }))}
          margin={{
            top: 5,
            right: 30,
            left: 20,
            bottom: 5,
          }}
        >
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis 
            dataKey="index" 
            domain={[0, chartData.miles_histogram.length + 5]}
            type = "number"
          />
          <YAxis />
          <Tooltip content={CustomTooltip}/>
          {/* <Legend /> */}
          <Bar dataKey="value" fill="#8884d8" />
        </BarChart>
    </>
  );
};

const CustomTooltip: ContentType<ValueType, NameType> = ({ active, payload, label }) => {
  if (active && payload && payload.length) {
    return (
      <div className="bg-white p-2 border border-gray-300 rounded shadow">
        <p><strong>Miles:</strong> {label}</p>
        <p><strong># Rides:</strong> {payload[0].value}</p>
      </div>
    );
  }

  return null;
};

