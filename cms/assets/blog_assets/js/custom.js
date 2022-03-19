jQuery.noConflict()(function($) {
    "use strict";
    if ($("div").is("#clock-seconds")) {
        (function() {
            var clockSeconds = document.getElementById("clock-seconds"),
                clockMinutes = document.getElementById("clock-minutes"),
                clockHours = document.getElementById("clock-hours");

            function getTime() {

                var date = new Date(),
                    seconds = date.getSeconds(),
                    minutes = date.getMinutes(),
                    hours = date.getHours(),

                    degSeconds = seconds * 360 / 60,
                    degMinutes = (minutes + seconds / 60) * 360 / 60,
                    degHours = (hours + minutes / 60 + seconds / 60 / 60) * 360 / 12;

                clockSeconds.setAttribute("style", "-webkit-transform: rotate(" + degSeconds + "deg); -moz-transform: rotate(" + degSeconds + "deg); -ms-transform: rotate(" + degSeconds + "deg); -o-transform: rotate(" + degSeconds + "deg); transform: rotate(" + degSeconds + "deg);");
                clockMinutes.setAttribute("style", "-webkit-transform: rotate(" + degMinutes + "deg); -moz-transform: rotate(" + degMinutes + "deg); -ms-transform: rotate(" + degMinutes + "deg); -o-transform: rotate(" + degMinutes + "deg); transform: rotate(" + degMinutes + "deg);");
                clockHours.setAttribute("style", "-webkit-transform: rotate(" + degHours + "deg); -moz-transform: rotate(" + degHours + "deg); -ms-transform: rotate(" + degHours + "deg); -o-transform: rotate(" + degHours + "deg); transform: rotate(" + degHours + "deg);");
            }

            setInterval(getTime, 1000);
            getTime();

        }());
    };

    // PORTFOLIO FILTERING - ISOTOPE
    //**********************************
    if ($("div").is(".oi_port_container")) {
        var $container = $('.oi_port_container');

        if ($container.length) {
            $container.waitForImages(function() {

                // initialize isotope
                $container.isotope({
                    itemSelector: '.oi_strange_portfolio_item',
                    layoutMode: 'masonry'
                });

                $('#filters a:first-child').addClass('filter_current');
                // filter items when filter link is clicked

                $("a", "#filters").on("click", function(e) {
                    var selector = $(this).attr('data-filter');
                    $container.isotope({
                        filter: selector
                    });
                    $(this).removeClass('filter_button').addClass('filter_button filter_current').siblings().removeClass('filter_button filter_current').addClass('filter_button');

                    return false;
                });
            }, null, true);
        }
    };


    if ($("div").is("#map")) {
        $("#map").gmap3({
            marker: {
                // address:"93 Worth St, New York, NY",
                address: "7th Ave, New York, NY",
                options: {
                    icon: "assets/css/img/marker.png"
                }
            },
            map: {
                options: {
                    styles: [{
                        stylers: [{
                            "saturation": -100
                        }, {
                            "lightness": 0
                        }, {
                            "gamma": 0.5
                        }]
                    }, ],
                    zoom: 13,
                    scrollwheel: false,
                    draggable: true
                }
            }
        });
    }

    if ($("div").is(".f_slider")) {
        $('.f_slider').flexslider({
            prevText: "", //String: Set the text for the "previous" directionNav item
            nextText: "",
            animation: "fade",
            useCSS: false,
            controlNav: false,
            animationLoop: true,
            slideshow: true,
            slideshowSpeed: 3000,
            pauseOnHover: true,
            start: function(slider) {
                slider.removeClass('oi_flex_loading');
            }
        });
    }
	
    if ($('body').width() > 640) {
        $(window).load(function() {
            if (($("body").height() - $(window).height()) > 300) {
                var stickyNavTop = $('.oi_head_holder').offset().top + $(".oi_head_holder .row").outerHeight();
                $(window).scroll(function() {
                    if ($(this).scrollTop() > stickyNavTop) {
                        $('.oi_st_menu_holder').fadeIn('fast');
                    } else {
                        $('.oi_st_menu_holder').fadeOut('fast');
                    }
                });
            };
        });
    };

    $(".oi_vc_clock", ".row").on("click", function(e) {
        $('#clockmodal').modal('toggle');
    });


    $(window).load(function() {
        $('#blog_snipet_slider').flexslider({
            animation: "fade",
            controlNav: true,
            directionNav: false,
            prevText: "Previous",
            nextText: "Next"
        });
    });


    if ($("div").is(".oi_head_holder")) {
        var stickyNavTopp = '';
        var stickyNavTop = $('.oi_head_holder').offset().top;
        stickyNavTopp = stickyNavTop + $('.oi_head_holder').outerHeight();
        $(window).scroll(function() {
            if ($(this).scrollTop() > stickyNavTopp) {
                $('.oi_head_holder').addClass('oi_sticky');
            } else {
                $('.oi_head_holder').removeClass('oi_sticky');
            }
        });
    }
    $(".oi_xs_menu", ".io_xs").on("click", function(e) {
        $('.oi_header_menu').toggleClass('oi_v_menu');
    });
	
	$(document).ready(function ()
			{ // after loading the DOM
				$("#ajax-contact-form").submit(function ()
				{
					// this points to our form
					var str = $(this).serialize(); // Serialize the data for the POST-request
					var result = '';
					$.ajax(
					{

						type: "POST",
						url: 'contact.php',
						data: str,
						success: function (msg)
						{	
								if (msg == 'OK'){
									result = '<div class="alert alert-info">Message was sent to website administrator, thank you!</div>';
									$("#fields").hide();
								}else{
									result = msg;
								}
								$("#note").html(result);
						
						}
					});
					return false;
				});
			});


});