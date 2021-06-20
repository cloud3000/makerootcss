# makerootcss

A simple utility to create css --color variables, and store them in a :root CSS pseudo-class { properties }.

Currently makerootcss will generate over 100 color variables, base on one initial color.

> EXAMPLE, typical scenario

```bash
./makerootcss --red 60 --green 120 --blue 205  > ./css/root.css
cat ./css/[YOUR CSS FILE] >> ./css/root.css
rm ./css/[YOUR CSS FILE]
mv ./css/root.css ./css/[YOUR CSS FILE]
```

For example:
Based on the RGB 60,120,205 hex: #3c78cd; will be the primary color.

makerootcss will generate a compliment color, I call this the secondary color, which will be 30 degrees forward on the color wheel, hex: #473bce;

The contrast or tertiary color will be on the opposite side of the color wheel, from between the primary and secondary. So from the primary +15, and then +180, so it's +195 degrees from the primary. This will be #cbb234; Think of this like placing a "Y" of the color wheel. From left to right the first point of the "Y" is the primary color, the second point is the secondary color, and the bottom point of the "Y" is the tertiary color.

The primary hue is 215 degrees, adding 195 to 215 gives you 50 because the color wheel is 360 degrees.

Additionally it will create an alarm color, a caution color, and a good color. These colors are for messaging, like BootStrap's success, warning, and danger. The will start off #ff0000; #ffff00; and #00ff00; however their saturation will be modified to match the saturation of the other colors.

Next it will also generate a gray color. Every good design need some gray.

Next it will also generate 12 shades of each color, from dark to light, 1 being the darkest and 12 being the lightest. The sixth color is the original shade, and there will also be three addition variables of the sixth shade, with opacities of .1, .2, and .3

The idea is that in your css you primarily use the sixth shade, for example; A button hover, you would set the button background to the primary06, and on hover set it to primary04, or 03 for more of an effect, remember 1 is the darkest, and 12 is the lightest. Also using the lighter colors give you more of a pastel look.

Gray is the exception, with 9 shades from black to white, and 10 different opacities.

Lastly; there is a config file rootcss.json where you can set the variable names, and the base color of the messaging colors. I use traffic-light colors, you may choose others.

> TODO's:
> Fix the hsl2rgb function. Currently the color variable are written out using HSL values, because hsl2rgb is not working. When I fix it I will output colors as hex rgb ea. #de45a3;

> Add geometry variables for height, width, size, padding, and margins for all the various types of selectors.

> Add the yspread logic; If you do not like the way this algorithm picks matching colors, there could be a property in the config file called yspread. It will allow you to change the distance between the points of the "Y" on the color wheel. It's set to 30 like I explained above. However, say you change it to 60. Now the secondary color will be +60 degrees from the primary, the the tertiary color will be +110 degrees from the primary. You can play with the yspread value and see what works best.
