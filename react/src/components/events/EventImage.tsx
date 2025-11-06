import React, { memo } from "react";
import { Image } from "../../models/UserModels";
import Slider from "react-slick";

type EventImagesProps = {
  images: Image[];
};

export const EventImages = memo<EventImagesProps>(({ images }) => {
  console.log({ images });
  if (images.length <= 0) return null;

  const isLocalhost = window.location.hostname === "localhost";

  var settings = {
    dots: true,
    speed: 500,
    slidesToShow: 1,
    slidesToScroll: 1,
    autoplay: true,
    infinite: images.length > 1,
  };

  return (
    <Slider {...settings}>
      {images.map((image) => (
        <div
          key={image.id}
          className="flex justify-center items-center h-60 p-4"
        >
          <img
            src={
              image.link.startsWith("https://")
                ? image.link
                : isLocalhost
                  ? `http://localhost:3000${image.link}`
                  : `/api${image.link}`
            }
            className="h-full w-auto object-cover rounded shadow-lg m-auto"
          />
        </div>
      ))}
    </Slider>
  );
});
