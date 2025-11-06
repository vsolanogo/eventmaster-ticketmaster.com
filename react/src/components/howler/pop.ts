import { Howl } from "howler";
import pop4 from "../../sounds/pop4.wav";

export const pop = new Howl({
  src: [pop4],
  autoplay: false,
  loop: false,
});
