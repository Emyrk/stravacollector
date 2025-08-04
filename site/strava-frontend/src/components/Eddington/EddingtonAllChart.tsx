import { FC, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { getAllAthleteEddingtons, getAthleteEddington } from "../../api/rest";
import { Eddington } from "../../api/typesGenerated";
import { useQuery } from "@tanstack/react-query";
import { Loading } from "../Loading/Loading";
import { ErrorBox } from "../ErrorBox/ErrorBox";
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Line, ResponsiveContainer, Label, ReferenceDot, Brush, ReferenceLine } from 'recharts';
import { ContentType } from "recharts/types/component/Tooltip";
import { NameType, ValueType } from "recharts/types/component/DefaultTooltipContent";
import { useTheme } from "@chakra-ui/react";
import { BrushStartEndIndex } from "recharts/types/context/brushUpdateContext";


export const EddingtonAllChart: FC<{}> = ({}) => {
  const [ zoomRange, setZoomRange ] = useState<[number, number] | undefined>(undefined);
  const theme = useTheme();


  const queryKey = ["athletes", "eddington"];
    const {
      data: chartData,
      error: chartError,
      isLoading: chartLoading,
      isFetched: chartFetched,
    } = useQuery({
      queryKey,
      queryFn: () =>
        getAllAthleteEddingtons(),
      onSuccess: (data) => {
      
      },
      onError: (error) => {
        console.error("Error fetching athlete data:", error);
      }
    });

  if (
    (!chartData || chartLoading) && !chartError
  ) {
    return <Loading />;
  }

  if (
    chartError || !chartData
  ) {
    return <ErrorBox error="Error fetching all athlete eddington data." detail={chartError} />;
  }
    

  const highest = Math.max.apply(Math, chartData.map((value) => value.current_eddington))


  const barData = chartData.reduce<number[]>((prev, value) => {
    prev[value.current_eddington-1]++
    return prev
  }, new Array<number>(highest)).
  map((value, index) => ({
    index:index+1, value,
  }))

  // const zoomedDomain =
  // zoomRange !== undefined
  //   ? [zoomRange[0], zoomRange[1]]
  //   : [1, chartData.miles_histogram.length];

  return (
    <>
    <div style={{ width: '100%', height: 300, position: 'relative' }}>
      <div style={{
        position: 'absolute',
        top: 10,
        right: 10,
        fontWeight: 'bold',
        fontSize: '20px',
        zIndex: 10,
        background: 'rgba(0,0 , 0, 0.6)',
        padding: '5px',
      }}>
        {/* Eddington Number = <span style={{color: theme.colors.brand.stravaOrange}}>{chartData.current_eddington}</span> */}
      </div>


    <ResponsiveContainer width="100%" height={300}>
        <BarChart
          data={barData}
          // margin={{
          //   top: 5,
          //   right: 30,
          //   left: 20,
          //   bottom: 5,
          // }}
        >
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis 
            dataKey="index" 
            // type = "number"
            // ticks={Array.from({ length: Math.floor(maxX / 25) + 1 }, (_, i) => i * 25)}
          />
          <YAxis 
          />
          <Tooltip content={CustomTooltip}/>
          {/* <Legend /> */}
          {/* Before "#8884d8" */}
          <Bar dataKey="value" fill="#8884d8" />
          <Line
            type="linear"
            dataKey="index"
            stroke="rgba(255, 0, 0, 0.6)"
            dot={false}
            isAnimationActive={false}
          />
          {/*  Other option */}
          {/* <ReferenceDot x={200} y={200} r={0} fill="none">
            <Label value="Eddington Number" position="top" offset={10} />
          </ReferenceDot> */}
          {/* <ReferenceLine
            x={chartData.current_eddington}
            stroke={theme.colors.brand.stravaOrange}
            strokeDasharray="6 6"
            label={{
              position: 'insideTopLeft',
              value: `${chartData.current_eddington}`,
              fill: theme.colors.brand.stravaOrange,
              fontSize: 12
            }}
          /> */}
          <Brush 
            dataKey="index" 
            height={30} 
            stroke={theme.colors.brand.stravaOrange} 
            startIndex={0}
            // endIndex={Math.min(chartData.current_eddington*2, chartData.miles_histogram.length)}
            onChange={(range) => {
              // props.startIndex
              // setZoomRange([range.startIndex, range.endIndex])
              // console.log("Brush changed:", range);
            }}
          />
        </BarChart>
      </ResponsiveContainer>
      </div>
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


