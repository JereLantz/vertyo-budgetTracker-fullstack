//@ts-check

/** update the visible saldo when new transaction is added to the page
 * @param {number} id 
 */
function udateSaldoSum(id){
    //TODO:
}

/** calculate the sum on page load, if the database contains existing transactions */
function calculateInitialSaldoSum(){
    let sum = 0
    const sumDisplay = document.getElementById("totalDispl")

    const transactions = document.querySelectorAll(".amountDisplay")

    for(const ta of transactions){
        let transaText = ta.textContent
        let parsedText = transaText?.replace("€", "")
        sum += Number(parsedText)
    }

    if (sumDisplay != null){
        sumDisplay.textContent = `${sum}€`
    }
}
