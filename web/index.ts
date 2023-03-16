async function getText(url: string) {
    const response = await fetch(url);
    return await response.text();
}


async function main() {

}

main();