/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("zpmsid5xlf5361s")

  // remove
  collection.schema.removeField("vmrrtaqw")

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("zpmsid5xlf5361s")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "vmrrtaqw",
    "name": "registry",
    "type": "select",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "maxSelect": 1,
      "values": [
        "github",
        "docker",
        "harbor"
      ]
    }
  }))

  return dao.saveCollection(collection)
})
