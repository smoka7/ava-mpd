import ColorThief from "colorthief";

const colorThief = new ColorThief();

export const colorSchemes = {
    first: [[51, 54, 82], [250, 208, 44], [189, 239, 83], [144, 173, 198], [213, 227, 172]],
    second: [[73, 88, 103], [87, 115, 153], [254, 95, 85,], [1189, 213, 234], [229, 229, 243]],
    third: [[11, 72, 107], [254, 67, 101], [252, 157, 154], [200, 200, 169], [249, 205, 173]],
    fourth: [[37, 41, 52], [201, 126, 24], [235, 158, 5], [227, 185, 107], [227, 185, 147]]
}

export function getRGB(color) {
    return "rgb(" + color[0] + "," + color[1] + "," + color[2] + ")";
}

function luminance(color) {
    let lum = [];
    for (let i = 0; i < 3; i++) {
        let chanel = color[i] / 255;
        lum[i] = (chanel <= 0.03928) ?
            chanel / 12.92 :
            Math.pow(((chanel + 0.055) / 1.055), 2.4);
    }
    return 0.2126 * lum[0] + 0.7152 * lum[1] + 0.0722 * lum[2];
}

function invert(color) {
    //https://24ways.org/2010/calculating-color-contrast
    var yiq = ((color[0] * 299) + (color[1] * 587) + (color[2] * 114)) / 1000;
    return (yiq >= 128) ? 'black' : 'white';
}

export function writeColors(colors = colorSchemes.first) {
    let style =
        ":root{--primary-color:" +
        getRGB(colors[0]) +
        ";--secondary-color:" +
        getRGB(colors[1]) +
        ";--accent-color:" +
        getRGB(colors[2]) +
        ";--lighter-color:" +
        getRGB(colors[3]) +
        ";--lightest-color:" +
        getRGB(colors[4]) +
        ";--foreground-color:" +
        invert(colors[0]) +
        ";}";
    document.getElementById("style").innerText = style;
};

function generateColors() {
    const img = document.getElementById("albumArt");
    const createPalette = () => {
        let colors = colorThief.getPalette(img, 5, 5);
        colors = colors.sort((a, b) => {
            return luminance(a) - luminance(b);
        });
        writeColors(colors);
    }
    if (img.completed) {
        createPalette();
        return;
    }
    img.addEventListener("load", () => {
        createPalette();
    });
};

export function setColorScheme() {
    let scheme = localStorage.getItem("colorScheme");
    switch (scheme) {
        case "auto":
            generateColors();
            break;
        case "second":
            writeColors(colorSchemes.second)
            break;
        case "third":
            writeColors(colorSchemes.third)
            break;
        case "fourth":
            writeColors(colorSchemes.fourth)
            break;
        default:
            writeColors(colorSchemes.first)
    }
};