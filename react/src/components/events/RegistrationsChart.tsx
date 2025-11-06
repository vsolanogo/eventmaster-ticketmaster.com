import React, { useEffect, useState } from "react";
import { Bar } from "react-chartjs-2";
import {
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Chart,
  BarElement,
} from "chart.js";

Chart.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement
);

import axios from "axios";

interface RegistrationData {
  date: string;
  count: number;
}

interface EventParticipantsChartProps {
  eventId: string;
}

const EventParticipantsChart: React.FC<EventParticipantsChartProps> = ({
  eventId,
}) => {
  const [registrationData, setRegistrationData] = useState<RegistrationData[]>(
    []
  );

  useEffect(() => {
    const fetchRegistrationData = async () => {
      try {
        const response = await axios.get(
          `/api/participant/event/${eventId}/registrations-per-day`
        );
        setRegistrationData(response.data);
      } catch (error) {
        console.error("Error fetching registration data:", error);
      }
    };

    fetchRegistrationData();
  }, [eventId]);

  const dates = registrationData.map((data) => new Date(data.date));
  const counts = registrationData.map((data) => data.count);

  const data = {
    labels: dates.map((date) => date.toLocaleDateString()),
    datasets: [
      {
        label: "Registrations per Day",
        data: counts,
        fill: false,
        backgroundColor: "rgba(75,192,192,0.2)",
        borderColor: "rgba(75,192,192,1)",
        borderWidth: 1,
      },
    ],
  };

  return (
    <div>
      <h3 className="text-4xl font-extrabold text-indigo-900 mb-6 border-b-2 border-indigo-500 pb-2 mt-8 mb-6">
        Registrations per Day
      </h3>
      {/* <Line data={data} /> */}
      <Bar data={data} />
    </div>
  );
};

export default EventParticipantsChart;
