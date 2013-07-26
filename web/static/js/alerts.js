var silenceRow = null;
var silenceId = null;

function clearSilenceLabels() {
  $("#silence_filters_table").empty();
}

function addSilenceLabel(label, value) {
  if (!label) {
    label = "";
  }
  if (!value) {
    value = "";
  }
  $("#silence_filters_table").append(
        '<tr>' +
        '  <td><input class="input-large" name="silence_filter_label[]" type="text" placeholder="label regex" value="' + label + '" required></td>' +
        '  <td><input class="input-large" name="silence_filter_value[]" type="text" placeholder="value regex" value="' + value + '" required></td>' +
        '  <td><button type="button" class="btn del_label_button"><i class="icon-minus"></i></button></td>' +
        '</tr>');
  bindDelLabel();
}

function bindDelLabel() {
  $(".del_label_button").unbind("click");
  $(".del_label_button").click(function() {
    $(this).parents("tr").remove();
  });
}

function silenceJsonFromForm() {
  var filters = {};
  var labels = $('input[name="silence_filter_label[]"]');
  var values = $('input[name="silence_filter_value[]"]');
  for (var i = 0; i < labels.length; i++) {
    filters[labels.get(i).value] = values.get(i).value;
  }

  var endsAt = 0;
  if ($("#silence_ends_at").val() != "") {
    var picker = $("#ends_at_datetimepicker").data("datetimepicker");
    endsAt = Math.round(picker.getLocalDate().getTime() / 1000);
  }

  return JSON.stringify({
    CreatedBy: $("#silence_created_by").val(),
    EndsAtSeconds: endsAt,
    Comment: $("#silence_comment").val(),
    Filters: filters
  });
}

function createSilence() {
  $.ajax({
    type: "POST",
    url: "/api/silences",
    data: silenceJsonFromForm(),
    dataType: "text",
    success: function(data, textStatus, jqXHR) {
      location.reload();
    },
    error: function(data, textStatus, jqXHR) {
      alert("Creating silence failed: " + textStatus);
    }
  });
}

function updateSilence() {
  $.ajax({
    type: "POST",
    url: "/api/silences/" + silenceId,
    data: silenceJsonFromForm(),
    dataType: "text",
    success: function(data, textStatus, jqXHR) {
      location.reload();
    },
    error: function(data, textStatus, jqXHR) {
      alert("Updating silence failed: " + textStatus);
    }
  });
}

function getSilence(silenceId, successFn) {
  $.ajax({
    type: "GET",
    url: "/api/silences/" + silenceId,
    async: false,
    success: successFn,
    error: function(data, textStatus, jqXHR) {
      alert("Creating silence failed: " + textStatus);
    }
  });
}

function deleteSilence(silenceId, silenceRow) {
  $.ajax({
    type: "DELETE",
    url: "/api/silences/" + silenceId,
    success: function(data, textStatus, jqXHR) {
      silenceRow.remove();
      $("#del_silence_modal").modal("hide");
    },
    error: function(data, textStatus, jqXHR) {
      alert("Removing silence failed: " + textStatus);
    }
  });
}

function initNewSilence() {
  silenceId = null;
  $("#edit_silence_header, #edit_silence_btn").html("Create Silence");
  $("#edit_silence_form")[0].reset();
}

function populateSilenceLabels(form) {
  var labels = form.children('input[name="label[]"]');
  var values = form.children('input[name="value[]"]');
  for (var i = 0; i < labels.length; i++) {
    addSilenceLabel(labels.get(i).value, values.get(i).value);
  }
}

function init() {
  $.ajaxSetup({
    cache: false
  });

  $("#ends_at_datetimepicker").datetimepicker({
    language: "en",
    pickSeconds: false
  });

  $("#edit_silence_modal").on("shown", function(e) {
    $("#silence_created_by").focus();
  });

  $("#edit_silence_modal").on("hidden", function(e) {
    clearSilenceLabels();
  });

  // This is the "Silences" page button to open the dialog for creating a new
  // silence.
  $("#new_silence_btn").click(function() {
    initNewSilence();
  });

  // These are the "Alerts" page buttons to open the dialog for creating a new
  // silence for the alert / alert instance.
  $(".add_silence_btn").click(function() {
    initNewSilence();
    populateSilenceLabels($(this).parents("form"));
  });

  $("#add_filter_button").click(function() {
    addSilenceLabel("", "");
  });

  $("#edit_silence_form").submit(function() {
    if (silenceId != null) {
      updateSilence();
    } else {
      createSilence();
    }
    return false;
  });

  $(".edit_silence_btn").click(function() {
    $("#edit_silence_header, #edit_silence_btn").html("Update Silence");

    silenceRow = $(this).parents("tr");
    silenceId = silenceRow.find("input[name='silence_id']").val();
    $("#edit_silence_form input[name='silence_id']").val(silenceId);
    getSilence(silenceId, function(silence) {
      var picker = $("#ends_at_datetimepicker").data('datetimepicker');
      var endsAt = new Date(silence.EndsAtSeconds * 1000);
      picker.setLocalDate(endsAt);

      $("#silence_created_by").val(silence.CreatedBy);
      $("#silence_comment").val(silence.Comment);
      for (var f in silence.Filters) {
        addSilenceLabel(f, silence.Filters[f]);
      }
    });
  });

  // When clicking the "Remove Silence" button in the Silences table, save the
  // table row and silence ID to global variables (ugh) so they can be used
  // from the modal dialog to remove that silence.
  $(".del_silence_modal_btn").click(function() {
    silenceRow = $(this).parents("tr");
    silenceId = silenceRow.find("input[name='silence_id']").val();
  });

  // Deletion confirmation button action.
  $(".del_silence_btn").click(function() {
    deleteSilence(silenceId, silenceRow);
  });
}

$(init);
