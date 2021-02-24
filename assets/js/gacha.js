var table;

$(document).ready(function() {
    function BtnLoading(elem) {
        $(elem).attr("data-original-text", $(elem).html());
        $(elem).prop("disabled", true);
        $(elem).html('<i class="spinner-border spinner-border-sm"></i> Loading...');
    }

    function BtnReset(elem) {
        $(elem).prop("disabled", false);
        $(elem).html($(elem).attr("data-original-text"));
    }

    $('.gachaLogURLSubmit').click(function(e) {
        e.preventDefault();

        if (table) {
            table.destroy();
        }

        var elem = $(e.currentTarget);

        BtnLoading(elem);

        var logURL = $('#inputGachaLogURL').val();
        processGachaLog(logURL);
    });
});

function resetButton() {
    // reset button
    $('.gachaLogURLSubmit').prop("disabled", false);
    $('.gachaLogURLSubmit').html("Submit");
}

function processGachaLog(logURL) {
    var url = "/gacha/process/history";

    var request = $.ajax({
        url: url,
        method: "POST",
        data: JSON.stringify({
            logURL: btoa(logURL),
        }),
        dataType: "json",
        contentType: "application/json",
    });

    request.done(processGachaResult);
    request.catch(function (err) {
        alert(err);
        resetButton()
    });
}

function processGachaResult(r) {
    if (!r.data) {
        alert("no data");
        resetButton();
        return
    }

    var noviceChartCanvas = $('#noviceChart').get(0).getContext('2d');
    var noviceWishes = r.data["100"];
    mapDataToChart("Novice Wishes", noviceWishes, noviceChartCanvas);


    var permanentChartCanvas = $('#permanentChart').get(0).getContext('2d');
    var permanentWishes = r.data["200"];
    mapDataToChart("Permanent Wishes", permanentWishes, permanentChartCanvas);


    var charEventChartCanvas = $('#charEventChart').get(0).getContext('2d');
    var charEventWishes = r.data["301"];
    mapDataToChart("Character Event Wishes", charEventWishes, charEventChartCanvas);


    var weaponChartCanvas = $('#weaponChart').get(0).getContext('2d');
    var weaponWishes = r.data["302"];
    mapDataToChart("Weapon Wishes", weaponWishes, weaponChartCanvas);

    table = $('#wishData').DataTable( {
        data: dataToTable(r.data),
        columns: [
            {
                title: "Wish Type",
                data: "gacha_type",
            },
            {
                title: "Name",
                data: "name",
            },
            {
                title: "Type",
                data: "item_type",
            },
            {
                title: "Rank",
                data: "rank_type",
            },
            {
                title: "Time",
                data: "time",
            },
        ]
    } );

    $('#collapseResult').click();
    $('#collapseData').click();

    resetButton();
}

var gachaType = {
    "100": "Novice Wishes",
    "200": "Permanent Wishes",
    "301": "Character Event Wishes",
    "302": "Weapon Wishes",
}

function dataToTable(data) {
    var combined = [
        ...data["100"],
        ...data["200"],
        ...data["301"],
        ...data["302"],
    ];

    var result = combined.map(item => ({
        gacha_type: gachaType[item.gacha_type],
        name: item.name,
        item_type: item.item_type,
        rank_type: item.rank_type,
        time: item.time,
    }))

    return result;
}

function mapDataToChart(title, data, chart) {
    if (!data) {
        return
    }
    var wishData = {
        labels: [
            '3* Weapon',
            '4* Weapon',
            '5* Weapon',
            '4* Character',
            '5* Character',
        ],
        datasets: [
            {
                data: [data.filter(function(wish){
                    return wish.rank_type === "3" && wish.item_type === "Weapon";
                }).length,data.filter(function(wish){
                    return wish.rank_type === "4" && wish.item_type === "Weapon";
                }).length,data.filter(function(wish){
                    return wish.rank_type === "5" && wish.item_type === "Weapon";
                }).length,data.filter(function(wish){
                    return wish.rank_type === "4" && wish.item_type === "Character";
                }).length,data.filter(function(wish){
                    return wish.rank_type === "5" && wish.item_type === "Character";
                }).length],
                backgroundColor : [
                    '#3fa7e3',
                    '#8632e0',
                    '#fab032',
                    '#b28dd9',
                    '#e0ae5e'
                ],
            }
        ]
    }

    var options = {
        maintainAspectRatio : false,
        responsive : true,
        title: {
            display: true,
            text: title,
        }
    }

    var pieChart = new Chart(chart, {
        type: 'pie',
        data: wishData,
        options: options
    });
}