using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;
using System;

namespace BLOG.Models
{
    public class Comment
    {
        [BsonId]
        [BsonRepresentation(BsonType.ObjectId)]
        public string? Id { get; set; }

        [BsonElement("PostId")]
        public string PostId { get; set; }  // ID blog posta na koji komentar ide

        [BsonElement("AuthorName")]
        public string AuthorName { get; set; } //TO DO: staviti samo id usera i ta se preko drugih servisa dobavlja ostale info

        [BsonElement("AuthorEmail")]
        public string AuthorEmail { get; set; }

        [BsonElement("Text")]
        public string Text { get; set; }

        [BsonElement("CreatedAt")]
        public DateTime CreatedAt { get; set; } = DateTime.UtcNow;

        [BsonElement("UpdatedAt")]
        public DateTime UpdatedAt { get; set; } = DateTime.UtcNow;
    }
}
